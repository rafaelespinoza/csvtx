package task

import (
	"fmt"
	"os"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func MintToYNAB(p Params) error {
	if p.Outdir == "" {
		if outdir, err := os.MkdirTemp("", "csvtx_*"); err != nil {
			return err
		} else {
			fmt.Fprintf(p.LogDest, "files will be written to tempdir %q\n", outdir)
			p.Outdir = outdir
		}
	}
	if p.LogDest == nil {
		p.LogDest = os.Stderr
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

	return readParseMint(p.Infile, func(m *entity.Mint) error {
		if _, ok := filesByAccount[m.Account]; !ok {
			entry, err := initOutfile(m.Account, ynabHeaders, p.Outdir)
			if err != nil {
				return err
			}
			filesByAccount[m.Account] = &entry
		}
		csvWriter := filesByAccount[m.Account].w
		ynabTx := entity.YNAB{
			Date:     m.Date,
			Payee:    m.Description,
			Category: m.Category,
			Memo:     m.Notes,
			Amount:   m.Amount,
		}
		row := ynabTx.AsRow()
		return csvWriter.Write(row)
	})
}

var ynabHeaders = []string{
	"Date", "Payee", "Category", "Memo", "Outflow", "Inflow",
}
