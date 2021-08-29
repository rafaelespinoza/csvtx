package entity

import (
	"fmt"
	"time"
)

type MechanicsBank struct {
	Date         time.Time
	Description  string
	Memo         string
	AmountDebit  AmountSubunits
	AmountCredit AmountSubunits
	Balance      AmountSubunits
	CheckNumber  int
	Fees         AmountSubunits
}

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

type WellsFargo struct {
	Date        time.Time
	Amount      AmountSubunits
	Description string
}

type YNAB struct {
	Date     time.Time
	Payee    string
	Category string
	Memo     string
	Amount   AmountSubunits
}
