package task

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
	"github.com/rafaelespinoza/csvtx/internal/product/mint"
)

// Interpret the csv file at filepath as exported csv transactions from Mint.com
// and then invoke callback on all of the parsed rows. Typically, this argument
// would be WriteAcctFiles.
func ReadParseMint(filepath string, callback func([]mint.Transaction)) {
	csvReader := initCsvReader(filepath)
	parseCSV(csvReader, callback)
}

func initCsvReader(filepath string) *csv.Reader {
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("error opening file\n%v\n", err)
	}

	br := bufio.NewReader(file)
	return csv.NewReader(br)
}

func parseCSV(reader *csv.Reader, callback func([]mint.Transaction)) {
	// Ignore first line b/c it's usually headers, not data. It screws up
	// parsing of data
	_, firstLineErr := reader.Read()
	if firstLineErr != nil {
		log.Fatalln("error reading the first line!")
	}

	var transactions []mint.Transaction

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

func parseLine(line []string) mint.Transaction {
	// columns to ignore:
	// "Original Description" 	(index 2)
	// "Labels" 				(index 7)
	tt := line[4]

	mt := mint.Transaction{
		Date:            parseDate(line[0]),
		Description:     line[1],
		TransactionType: tt,
		Category:        line[5],
		Account:         line[6],
		Notes:           line[8],
	}

	isNegative, err := mt.Negative()

	if err != nil {
		log.Fatalf("error setting transaction type. %T, %v. %v\n", mt, mt, err)
	}

	mt.Amount = parseMoney(line[3], isNegative) // "1234.56"
	return mt
}

func parseMoney(cell string, isNegative bool) entity.AmountSubunits {
	s := fallbackStr(cell, "0.00")
	dc := strings.Split(s, ".")

	if len(dc) < 2 {
		dc = append(dc, "00")
	}

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

	return entity.AmountSubunits(a)
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
