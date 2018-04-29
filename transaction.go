package csvtx

import (
	"time"
)

type Transaction struct {
	Date            time.Time
	Description     string
	Amount          Amount // cents
	TransactionType string // [debit, credit]
	Category        string
	Account         string
	Notes           string

	// ignoring these columns from mint csv file:
	// "Original Description", "Labels"
}
