package csvtx

import (
	"fmt"
	"strings"
	"time"
)

type MintTransaction struct {
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

func (t MintTransaction) String() string {
	txType := strings.ToLower(t.TransactionType)

	return fmt.Sprintf(
		"{ Date: '%s', Description: '%s', Amount: %s, TransactionType: '%s', Category: '%s', Account: '%s', Notes: '%s' }",
		t.Date.Format("2006-01-02"),
		t.Description,
		t.Amount,
		txType,
		t.Category,
		t.Account,
		t.Notes,
	)
}
