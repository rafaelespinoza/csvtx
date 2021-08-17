package task

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/rafaelespinoza/csvtx/internal/product/mint"
	"github.com/rafaelespinoza/csvtx/internal/product/ynab"
)

// Creates a new csv file for every unique account found in transactions. Files
// are put into current working directory. Pass this function as a callback to
// ReadParseMint. Logs results to standard output.
func WriteAcctFiles(transactions []mint.Transaction) {
	uniqAcctNames := newAccountFilenames(&transactions)
	results := make([]string, 0, len(uniqAcctNames))

	// keep it simple, use multiple passes to write one file per account.
	for acct, file := range uniqAcctNames {
		result := <-writeFormatYnabAsync(transactions, acct, file)
		results = append(results, result)
	}

	log.Println(results)
}

func writeFormatYnabAsync(mints []mint.Transaction, targetAcctName, fileName string) <-chan string {
	c := make(chan string)

	go func() {
		c <- writeFormatYnab(mints, targetAcctName, fileName)
	}()

	return c
}

func writeFormatYnab(mints []mint.Transaction, targetAcctName, outputFile string) string {
	csvWriter := initCsvWriter(outputFile)

	defer csvWriter.dest.Close()

	csvWriter.writeHeaders(ynab.Headers)
	csvWriter.writeBody(mints, targetAcctName)

	log.Println("wrote file", outputFile)
	return outputFile
}

type csvWriter struct {
	dest   io.WriteCloser
	writer *csv.Writer
}

func initCsvWriter(filename string) csvWriter {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		log.Fatalln("error opening file")
	}

	return csvWriter{dest: file, writer: csv.NewWriter(file)}
}

func (w *csvWriter) writeHeaders(headers []string) {
	w.writer.Write(headers)
	w.writer.Flush()
}

func (w *csvWriter) writeBody(mints []mint.Transaction, targetAcctName string) {
	var ynab ynab.Transaction
	var row []string
	var txAcctType string

	for _, mintTx := range mints {
		txAcctType = mintTx.Account

		if txAcctType == targetAcctName {
			ynab = mintToYnab(mintTx)
			row = ynab.AsRow()
			w.writer.Write(row)
		}
	}

	w.writer.Flush()

	if err := w.writer.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func mintToYnab(mt mint.Transaction) ynab.Transaction {
	return ynab.Transaction{
		Date:     mt.Date,
		Payee:    mt.Description,
		Category: mt.Category,
		Memo:     mt.Notes,
		Amount:   mt.Amount,
	}
}
