package convert

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

type Params struct {
	Infile  string
	Outdir  string
	LogDest io.Writer
}

func (p *Params) init() error {
	if p.LogDest == nil {
		p.LogDest = os.Stderr
	}

	if p.Outdir == "" {
		if outdir, err := os.MkdirTemp("", "csvtx_*"); err != nil {
			return err
		} else {
			fmt.Fprintf(p.LogDest, "files will be written to tempdir %q\n", outdir)
			p.Outdir = outdir
		}
	}

	return nil
}

var ynabHeaders = []string{
	"Date", "Payee", "Category", "Memo", "Outflow", "Inflow",
}

func parseDate(inputDate string) (t time.Time, e error) {
	if t, e = time.Parse("1/02/2006", inputDate); e == nil {
		return
	}

	// try the other known format
	t, e = time.Parse("01/02/2006", inputDate)
	return
}

func parseMoney(cell string) (out entity.AmountSubunits, err error) {
	if cell == "" {
		return
	}
	tmp, err := strconv.ParseFloat(cell, 64)
	if err != nil {
		return
	}
	tmp *= 100
	tmp = math.Round(tmp)
	out = entity.AmountSubunits(tmp)
	return
}

type csvOut struct {
	f *os.File
	w *csv.Writer
}

func initOutfile(accountName string, headers []string, basedir string) (out csvOut, err error) {
	filename := strings.TrimSpace(accountName)
	filename = strings.Replace(filename, " ", "-", -1)
	filename = strings.ToLower(filename)
	filename = filepath.Join(basedir, filename+".csv")
	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		err = fmt.Errorf("could not create output file %q; %w", filename, err)
		return
	}

	w := csv.NewWriter(file)
	err = writeCSVHeaders(w, headers)
	if err != nil {
		err = fmt.Errorf("could not write csv headers; %w", err)
	}
	out = csvOut{f: file, w: w}
	return
}

func writeCSVHeaders(w *csv.Writer, headers []string) (err error) {
	if err = w.Write(headers); err != nil {
		return err
	}
	w.Flush()
	err = w.Error()
	return
}

func ynabAsRow(t entity.YNAB) []string {
	var outflow, inflow string
	if t.Amount < 0 {
		outflow = t.Amount.String()
	} else {
		inflow = t.Amount.String()
	}
	return []string{
		t.Date.Format("01/02/2006"),
		t.Payee,
		t.Category,
		t.Memo,
		outflow,
		inflow,
	}
}
