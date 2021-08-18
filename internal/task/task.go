package task

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
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

func parseDate(inputDate string) (t time.Time, e error) {
	if t, e = time.Parse("1/02/2006", inputDate); e == nil {
		return
	}

	// try the other known format
	t, e = time.Parse("01/02/2006", inputDate)
	return
}

func parseMoney(cell string, isNegative bool) (out entity.AmountSubunits, err error) {
	if cell == "" {
		return
	}
	tmp, err := strconv.ParseFloat(cell, 64)
	if err != nil {
		return
	}
	amt := tmp * 100
	if isNegative {
		amt *= -1
	}
	out = entity.AmountSubunits(amt)
	return
}

type csvOut struct {
	f *os.File
	w *csv.Writer
}

func initOutfile(accountType string, headers []string, basedir string) (out csvOut, err error) {
	filename := strings.TrimSpace(accountType)
	filename = strings.Replace(filename, " ", "-", -1)
	filename = strings.ToLower(filename)
	filename = path.Join(basedir, filename+".csv")
	file, err := os.Create(filename)
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
