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

type Transaction interface {
	AsRow() []string
}

type TransactionList interface {
	Export() []Transaction
}

func ReadParseMint(filepath string, callback func([]MintTransaction)) {
	csvFile, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("error opening file\n%v\n", err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	ParseCSV(reader, callback, true)
}

func ParseCSV(
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
	tt := parseString(line[4])

	mt := MintTransaction{
		Date:            parseDate(line[0]),
		Description:     parseString(line[1]),
		transactionType: tt,
		Category:        parseString(line[5]),
		Account:         parseString(line[6]),
		Notes:           parseString(line[8]),
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

func parseString(cell string) string {
	return cell
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
