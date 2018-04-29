package csvtx

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

func ReadParse(filepath string, callback func([]Transaction)) {
	csvFile, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("error opening file\n%v\n", err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	ParseCSV(reader, callback, true)
}
