package csvtx

import (
	"encoding/csv"
	"log"
	"os"
)

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
	csvWriter := initCsvWriter(outputFile)

	defer csvWriter.file.Close()

	csvWriter.writeHeaders(YnabHeader)
	csvWriter.writeBody(mints, targetAcctName)

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

func (w *CsvWriter) writeHeaders(headers []string) {
	w.writer.Write(headers)
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
