package task

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/rafaelespinoza/csvtx/internal/entity"
	"github.com/rafaelespinoza/csvtx/internal/product/mint"
)

// readParseMint interprets a CSV file at filepath as exported transactions from
// Mint.com and invokes onRow to handle a parsed CSV row.
func readParseMint(filepath string, callback func(*mint.Transaction) error) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(bufio.NewReader(file))
	return parseMintCSV(csvReader, callback)
}

func parseMintCSV(reader *csv.Reader, onRow func(*mint.Transaction) error) error {
	// Ignore first line b/c it's usually headers, not data. It screws up
	// parsing of data
	lineNumber := 1
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("could not read line %d; %w", lineNumber, err)
	}

	for {
		lineNumber++

		line, err := reader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("could not read line %d; %w", lineNumber, err)
		}

		tx, err := parseMintRow(line)
		if err != nil {
			return fmt.Errorf("could not parse line %d; %w", lineNumber, err)
		}

		if err = onRow(tx); err != nil {
			return fmt.Errorf("callback error line %d; %w", lineNumber, err)
		}
	}
}

func parseMintRow(in []string) (out *mint.Transaction, err error) {
	// columns to ignore:
	// "Original Description" 	(index 2)
	// "Labels" 				(index 7)
	tt := in[4]

	date, err := parseDate(in[0])
	if err != nil {
		return
	}

	mt := mint.Transaction{
		Date:            date,
		Description:     in[1],
		TransactionType: tt,
		Category:        in[5],
		Account:         in[6],
		Notes:           in[8],
	}

	isNegative, err := mt.Negative()
	if err != nil {
		return
	}

	amount, err := parseMoney(in[3], isNegative) // "1234.56"
	if err != nil {
		return
	}
	mt.Amount = amount
	out = &mt
	return
}
