package convert

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func VenmoToYNAB(p Params) error {
	if err := p.init(); err != nil {
		return err
	}

	const accountName = "venmo"
	output, err := initOutfile(accountName, ynabHeaders, p.Outdir)
	if err != nil {
		return err
	}
	defer func() {
		var err error
		outfile := output.f.Name()

		output.w.Flush()

		if err = output.w.Error(); err != nil {
			fmt.Fprintf(p.LogDest, "could not flush %q data to %q; %v\n", accountName, outfile, err)
		}
		if err = output.f.Close(); err != nil {
			fmt.Fprintf(p.LogDest, "could not close file %q; %v\n", outfile, err)
		}
		if err == nil {
			fmt.Fprintf(p.LogDest, "wrote %q file %q\n", accountName, outfile)
		}
	}()

	return readParseVenmo(p.Infile, func(m *entity.Venmo) error {
		payee := m.From
		if m.Amount < 0 && m.TransactionType == "Payment" {
			payee = m.To
		} else if m.Amount > 0 && m.TransactionType == "Charge" {
			payee = m.To
		}
		row := ynabAsRow(entity.YNAB{
			Date:   m.Datetime,
			Payee:  payee,
			Memo:   m.Note,
			Amount: m.Amount,
		})
		return output.w.Write(row)
	})
}

func readParseVenmo(pathToFile string, onRow func(*entity.Venmo) error) error {
	file, err := os.Open(filepath.Clean(pathToFile))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	csvReader := csv.NewReader(bufio.NewReader(file))
	const countNonDataRows = 4
	for i := 0; i < countNonDataRows; i++ {
		// Read past some header rows until the data starts.
		if _, err := csvReader.Read(); err != nil {
			return fmt.Errorf("could not read line %d; %w", i+1, err)
		}
	}

	lineNumber := countNonDataRows
	for {
		lineNumber++

		line, err := csvReader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("could not read line %d; %w", lineNumber, err)
		}

		tx, err := parseVenmoRow(line)
		if err == errNotTransaction {
			// skip over the data rows that only denote account balance.
			continue
		} else if err != nil {
			return fmt.Errorf("could not parse line %d; %w", lineNumber, err)
		}

		if err = onRow(tx); err != nil {
			return fmt.Errorf("onRow error line %d; %w", lineNumber, err)
		}
	}
}

var errNotTransaction = errors.New("not a transaction")

func parseVenmoRow(in []string) (out *entity.Venmo, err error) {
	if in[1] == "" { // Consider an ID value a pre-requisite for a transaction.
		err = errNotTransaction
		return
	}
	date, err := time.Parse("2006-01-02T15:04:05", in[2])
	if err != nil {
		return
	}

	m := strings.Replace(in[8], " ", "", -1)
	m = strings.Replace(m, "$", "", -1)
	m = strings.Replace(m, "+", "", -1)
	m = strings.Replace(m, "-", "", -1)
	amount, err := parseMoney(m) // "- $123.45", "+ $78.90"
	if err != nil {
		return
	}
	if strings.HasPrefix(in[8], "-") {
		amount *= -1
	}

	out = &entity.Venmo{
		Datetime:        date,
		TransactionType: in[3],
		Note:            in[5], // probably will have emoji
		From:            in[6],
		To:              in[7],
		Amount:          amount,
	}
	return
}
