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
	return MintTransaction{
		Date:            parseDate(line[0]),
		Description:     parseString(line[1]),
		Amount:          parseMoney(line[3]), // "1234.56"
		TransactionType: parseString(line[4]),
		Category:        parseString(line[5]),
		Account:         parseString(line[6]),
		Notes:           parseString(line[8]),
	}
}

func parseMoney(cell string) Amount {
	s := fallbackStr(cell, "0.00")
	dc := strings.Replace(s, ".", "", 1)
	m, e := strconv.ParseInt(dc, 0, 0)

	if e != nil {
		log.Fatalf("error parsing money: %v", cell)
		return 0
	}

	return Amount(int(m))
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
