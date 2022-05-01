package convert

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func MintToYNAB(p Params) error {
	if err := p.init(); err != nil {
		return err
	}

	filesByAccount := make(map[string]*csvOut)
	defer func() {
		for accountType, csv := range filesByAccount {
			var err error
			outfile := csv.f.Name()

			csv.w.Flush()

			if err = csv.w.Error(); err != nil {
				fmt.Fprintf(p.LogDest, "could not flush %q data to %q; %v\n", accountType, outfile, err)
			}
			if err = csv.f.Close(); err != nil {
				fmt.Fprintf(p.LogDest, "could not close file %q; %v\n", outfile, err)
			}
			if err == nil {
				fmt.Fprintf(p.LogDest, "wrote %q file %q\n", accountType, outfile)
			}
		}
	}()

	infile, err := openFile(p.Infile)
	if err != nil {
		return err
	}
	defer func() { _ = infile.Close() }()

	return readParseMintCSV(infile, func(m *entity.Mint) error {
		if _, ok := filesByAccount[m.Account]; !ok {
			entry, err := initOutfile(m.Account, ynabHeaders, p.Outdir)
			if err != nil {
				return err
			}
			filesByAccount[m.Account] = &entry
		}
		csvWriter := filesByAccount[m.Account].w
		row := ynabAsRow(entity.YNAB{
			Date:     m.Date,
			Payee:    m.Description,
			Category: m.Category,
			Memo:     m.Notes,
			Amount:   m.Amount,
		})
		return csvWriter.Write(row)
	})
}

func readParseMintCSV(r io.Reader, onRow func(*entity.Mint) error) error {
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

		tx, err := parseMintRow(line)
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

func parseMintRow(in []string) (out *entity.Mint, err error) {
	// Is this a header row?
	if in[0] == "Date" {
		err = errNotTransaction
		return
	}

	date, err := parseDate(in[0])
	if err != nil {
		return
	}

	mt := entity.Mint{
		Date:            date,
		Description:     in[1],
		TransactionType: in[4],
		Category:        in[5],
		Account:         in[6],
		Notes:           in[8],
	}

	amount, err := parseMoney(in[3]) // "1234.56"
	if err != nil {
		return
	}
	negative, err := mt.Negative()
	if err != nil {
		return
	} else if negative {
		amount *= -1
	}
	mt.Amount = amount
	out = &mt
	return
}
