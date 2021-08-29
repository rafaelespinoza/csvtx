package convert

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func MechanicsBankToYNAB(p Params) error {
	if err := p.init(); err != nil {
		return err
	}

	const accountName = "mechanicsbank"
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

	return readParseMechanicsBank(p.Infile, func(m *entity.MechanicsBank) error {
		var amount entity.AmountSubunits
		if m.AmountCredit == 0 && m.AmountDebit != 0 { // is it negative?
			amount = m.AmountDebit * -1
		} else {
			amount = m.AmountCredit
		}
		row := ynabAsRow(entity.YNAB{
			Date:   m.Date,
			Payee:  m.Description,
			Memo:   m.Memo,
			Amount: amount,
		})
		return output.w.Write(row)
	})
}

func readParseMechanicsBank(filepath string, onRow func(*entity.MechanicsBank) error) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(bufio.NewReader(file))
	// There are some metadata rows at the top, with varying numbers of columns.
	// Reading past them for now.
	csvReader.FieldsPerRecord = -1

	// The input CSV will either include or exclude a Balance field. So, all the
	// data rows in some files might contain 8 fields. But if you were to
	// download another CSV a few days later, all the data rows could have 9
	// fields. I've witnessed it going back and forth; it's seemingly random.
	//
	// For the immediate future, there isn't a use case to use the last few
	// columns anyways, so it doesn't matter; those values are ignored. Some
	// extra logic would be needed to handle this variation if those values at
	// the end are needed. See the tests and testdata for known examples.
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

		tx, err := parseMechanicsBankRow(line)
		if err != nil {
			return fmt.Errorf("could not parse line %d; %w", lineNumber, err)
		}

		if err = onRow(tx); err != nil {
			return fmt.Errorf("onRow error line %d; %w", lineNumber, err)
		}
	}
}

func parseMechanicsBankRow(in []string) (out *entity.MechanicsBank, err error) {
	var (
		date                      time.Time
		amountDebit, amountCredit entity.AmountSubunits
	)

	if date, err = parseDate(in[1]); err != nil {
		err = fmt.Errorf("col %d; %w", 1, err)
		return
	}

	if amountDebit, err = parseMoney(in[4], true); err != nil { // -1234.56
		err = fmt.Errorf("col %d; %w", 4, err)
		return
	}
	if amountCredit, err = parseMoney(in[5], false); err != nil { // 1234.56
		err = fmt.Errorf("col %d; %w", 5, err)
		return
	}

	out = &entity.MechanicsBank{
		Date:         date,
		Description:  in[2],
		AmountDebit:  amountDebit,
		AmountCredit: amountCredit,
	}
	return
}
