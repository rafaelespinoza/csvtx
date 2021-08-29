package convert

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func WellsFargoToYNAB(p Params) error {
	if err := p.init(); err != nil {
		return err
	}

	const accountName = "wellsfargo"
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

	return readParseWellsFargo(p.Infile, func(m *entity.WellsFargo) error {
		row := ynabAsRow(entity.YNAB{
			Date:   m.Date,
			Payee:  m.Description,
			Amount: m.Amount,
		})
		return output.w.Write(row)
	})
}

func readParseWellsFargo(filepath string, onRow func(*entity.WellsFargo) error) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(bufio.NewReader(file))

	var lineNumber int
	for {
		lineNumber++

		line, err := csvReader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("could not read line %d; %w", lineNumber, err)
		}

		tx, err := parseWellsFargoRow(line)
		if err != nil {
			return fmt.Errorf("could not parse line %d; %w", lineNumber, err)
		}

		if err = onRow(tx); err != nil {
			return fmt.Errorf("onRow error line %d; %w", lineNumber, err)
		}
	}
}

func parseWellsFargoRow(in []string) (out *entity.WellsFargo, err error) {
	date, err := parseDate(in[0])
	if err != nil {
		return
	}

	isNegative := strings.HasPrefix(in[1], "-")
	amount, err := parseMoney(in[1], isNegative) // "-1234.56", "78.90"
	if err != nil {
		return
	}

	out = &entity.WellsFargo{
		Date:        date,
		Description: in[4],
		Amount:      amount,
	}
	return
}
