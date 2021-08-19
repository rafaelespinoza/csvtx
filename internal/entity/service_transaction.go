package entity

import (
	"fmt"
	"time"
)

type Mint struct {
	Date            time.Time
	Description     string
	Category        string
	Account         string
	Notes           string
	Amount          AmountSubunits
	TransactionType string // [debit, credit]

	// ignoring these columns from mint csv file:
	// "Original Description", "Labels"
}

func (t Mint) Negative() (bool, error) {
	switch tt := t.TransactionType; tt {
	case "debit":
		return true, nil
	case "credit":
		return false, nil
	default:
		return false, fmt.Errorf("invalid transaction type %q", tt)
	}
}

type YNAB struct {
	Date     time.Time
	Payee    string
	Category string
	Memo     string
	Amount   AmountSubunits
}

func (t YNAB) AsRow() []string {
	var outflow, inflow string
	if t.Amount < 0 {
		outflow = t.Amount.String()
	} else {
		inflow = t.Amount.String()
	}
	return []string{
		t.Date.Format("01/02/2006"),
		t.Payee,
		t.Category,
		t.Memo,
		outflow,
		inflow,
	}
}
