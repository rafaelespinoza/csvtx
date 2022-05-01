package convert

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
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

	infile, err := openFile(p.Infile)
	if err != nil {
		return err
	}
	defer func() { _ = infile.Close() }()

	return readParseVenmoCSV(infile, func(m *entity.Venmo) error {
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

func ReadParseVenmo(r io.Reader, sortDateAsc bool) (out []*entity.Venmo, err error) {
	onRow := func(row *entity.Venmo) error {
		out = append(out, row)
		return nil
	}
	err = readParseVenmoCSV(r, onRow)
	if err != nil {
		return
	}

	sort.SliceStable(out, func(i, j int) bool {
		if sortDateAsc {
			return out[i].Datetime.Before(out[j].Datetime)
		}
		return out[j].Datetime.Before(out[i].Datetime)
	})

	return
}

func readParseVenmoCSV(r io.Reader, onRow func(*entity.Venmo) error) error {
	csvReader := csv.NewReader(bufio.NewReader(r))
	var lineNumber int

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
			continue
		} else if err != nil {
			return fmt.Errorf("could not parse line %d; %w", lineNumber, err)
		}

		if err = onRow(tx); err != nil {
			return fmt.Errorf("onRow error line %d; %w", lineNumber, err)
		}
	}
}

func parseVenmoRow(in []string) (out *entity.Venmo, err error) {
	// Safety check. Also, there are more columns, but don't need them now.
	if len(in) < 9 {
		err = errNotTransaction
		return
	}

	// Is this a metadata row? Applicable to first 2 lines of a file.
	if strings.HasPrefix(in[0], "Account") {
		err = errNotTransaction
		return
	}
	// Is this a header row? 3rd line of the file.
	if in[0] == "" && in[1] == "ID" {
		err = errNotTransaction
		return
	}
	// Is this the row that only denotes account balance? 4th line of the file.
	if in[1] == "" { // Consider an ID value a pre-requisite for a transaction
		err = errNotTransaction
		return
	}

	// Now parse transaction data.

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
