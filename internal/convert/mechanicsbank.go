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

	infile, err := openFile(p.Infile)
	if err != nil {
		return err
	}
	defer func() { _ = infile.Close() }()

	return readParseMechanicsBankCSV(infile, func(m *entity.MechanicsBank) error {
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

func ReadParseMechanicsBank(r io.Reader, sortDateAsc bool) (out []*entity.MechanicsBank, err error) {
	onRow := func(row *entity.MechanicsBank) error {
		out = append(out, row)
		return nil
	}
	err = readParseMechanicsBankCSV(r, onRow)
	if err != nil {
		return
	}

	sort.SliceStable(out, func(i, j int) bool {
		if sortDateAsc {
			return out[i].Date.Before(out[j].Date)
		}
		return out[j].Date.Before(out[i].Date)
	})
	return
}

func readParseMechanicsBankCSV(r io.Reader, onRow func(*entity.MechanicsBank) error) error {
	csvReader := csv.NewReader(bufio.NewReader(r))
	// There are some metadata rows at the top, with varying numbers of columns.
	csvReader.FieldsPerRecord = -1

	var lineNumber int

	for {
		lineNumber++

		line, err := csvReader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("could not read line %d; %w", lineNumber, err)
		}

		tx, err := parseMechanicsBankRow(line)
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

func parseMechanicsBankRow(in []string) (out *entity.MechanicsBank, err error) {
	// Is this a metadata row?
	if len(in) < 6 {
		err = errNotTransaction
		return
	}

	// Is this a header row?
	if in[1] == "Date" {
		err = errNotTransaction
		return
	} else if strings.Contains(in[4], "Amount") || strings.Contains(in[5], "Amount") {
		err = errNotTransaction
		return
	}

	// Now parse transaction data.
	//
	// For what it's worth, sometimes a Mechanics Bank CSV export will contain a
	// Balance field and sometimes it won't; it's not very consistent. See the
	// tests and testdata for known examples. There isn't a use case right now
	// for the last few columns anyways, so the values are discarded.

	var (
		date                      time.Time
		amountDebit, amountCredit entity.AmountSubunits
	)

	if date, err = parseDate(in[1]); err != nil {
		err = fmt.Errorf("col %d; %w", 1, err)
		return
	}

	if amountDebit, err = parseMoney(in[4]); err != nil { // -1234.56
		err = fmt.Errorf("col %d; %w", 4, err)
		return
	}
	amountDebit *= -1
	if amountCredit, err = parseMoney(in[5]); err != nil { // 1234.56
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
