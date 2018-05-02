package csvtx

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadParseMint(filepath string, callback func([]MintTransaction)) {
	csvReader := initCsvReader(filepath)
	parseCSV(csvReader, callback, true)
}

func initCsvReader(filepath string) *csv.Reader {
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("error opening file\n%v\n", err)
	}

	br := bufio.NewReader(file)
	return csv.NewReader(br)
}

func WriteAcctFiles(transactions []MintTransaction) {
	uniqAcctNames := newAccountNameMap(&transactions)
	results := make([]string, 0, len(uniqAcctNames))

	// keep it simple, use multiple passes to write one file per account.
	for acct, file := range uniqAcctNames {
		result := <-writeFormatYnabAsync(&transactions, acct, file)
		results = append(results, result)
	}

	log.Println(results)
}

func writeFormatYnabAsync(
	mints *[]MintTransaction,
	targetAcctName,
	fileName string,
) <-chan string {
	c := make(chan string)

	go func() {
		c <- writeFormatYnab(mints, targetAcctName, fileName)
	}()

	return c
}

func writeFormatYnab(
	mints *[]MintTransaction,
	targetAcctName,
	outputFile string,
) string {
	w := initCsvWriter(outputFile)

	defer w.file.Close()

	w.writeHeaders()
	w.writeBody(mints, targetAcctName)

	log.Printf("wrote file %s\n", outputFile)
	return outputFile
}

type CsvWriter struct {
	file   *os.File
	writer *csv.Writer
}

func initCsvWriter(filename string) CsvWriter {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		log.Fatalln("error opening file")
	}

	return CsvWriter{file, csv.NewWriter(file)}
}

func (w *CsvWriter) writeHeaders() {
	w.writer.Write(YnabHeader)
	w.writer.Flush()
}

func (w *CsvWriter) writeBody(mints *[]MintTransaction, targetAcctName string) {
	var ynab YnabTransaction
	var row []string
	var txAcctType string

	for _, mintTx := range *mints {
		txAcctType = mintTx.Account

		if txAcctType == targetAcctName {
			ynab = mintTx.asYnabTx()
			row = ynab.AsRow()
			w.writer.Write(row)
		}
	}

	w.writer.Flush()

	if err := w.writer.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func parseCSV(
	reader *csv.Reader,
	callback func([]MintTransaction),
	ignoreHeaders bool,
) {
	if ignoreHeaders {
		// ignore first line (it's usually headers, not data) and it screws up
		// parsing of data
		_, firstLineErr := reader.Read()
		if firstLineErr != nil {
			log.Fatalln("error reading the first line!")
		}
	}

	var transactions []MintTransaction

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		transactions = append(transactions, parseLine(line))
	}

	callback(transactions)
}

func parseLine(line []string) MintTransaction {
	// columns to ignore:
	// "Original Description" 	(index 2)
	// "Labels" 				(index 7)
	tt := line[4]

	mt := MintTransaction{
		Date:            parseDate(line[0]),
		Description:     line[1],
		TransactionType: tt,
		Category:        line[5],
		Account:         line[6],
		Notes:           line[8],
	}

	isNegative, err := mt.isNegative()

	if err != nil {
		log.Fatalf("error setting transaction type. %T, %v. %v\n", mt, mt, err)
	}

	mt.Amount = parseMoney(line[3], isNegative) // "1234.56"
	return mt
}

func parseMoney(cell string, isNegative bool) Amount {
	s := fallbackStr(cell, "0.00")
	dc := strings.Split(s, ".")
	d, c := dc[0], dc[1] // dollars, cents

	var m int64
	var e error

	if d == "0" {
		m, e = strconv.ParseInt(c, 0, 0)
	} else {
		m, e = strconv.ParseInt(strings.Join(dc, ""), 0, 0)
	}

	if e != nil {
		log.Fatalf("error parsing money. cell: %v, dc: %v, m: %d", cell, dc, m)
		return 0
	}

	var a int

	if isNegative {
		a = int(m * -1)
	} else {
		a = int(m)
	}

	return Amount(a)
}

func parseDate(inputDate string) time.Time {
	var t time.Time
	var e error

	t, e = time.Parse("1/02/2006", inputDate)

	if e != nil {
		// try the other known format
		t, e = time.Parse("01/02/2006", inputDate)

		if e != nil {
			// now we're out of ideas
			log.Fatalf("could not parse date: %s\n", inputDate)
			return time.Time{}
		}
	}

	return t
}

func fallbackStr(input, fallback string) string {
	if len(input) < 1 {
		return fallback
	} else {
		return input
	}
}
