package csvtx

import (
	"fmt"
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

type MintToYnab []MintTransaction

func InitMintToYnab(mt *[]MintTransaction) MintToYnab {
	n := len(*mt)
	mty := make(MintToYnab, n, n)

	for i, t := range *mt {
		mty[i] = t
	}

	return mty
}

func (mty MintToYnab) Export() []YnabTransaction {
	n := len(mty)
	yt := make([]YnabTransaction, n, n)

	for i, t := range mty {
		yt[i] = t.asYnabTx()
	}

	return yt
}

func (mt MintTransaction) asYnabTx() YnabTransaction {
	return YnabTransaction{
		Date:     mt.Date, // hope this is not a reference
		CheckNum: 0,       // TODO: extract from t.Notes? t.Description?
		Payee:    mt.Description,
		Memo:     mt.Notes,
		Amount:   mt.Amount, // hopefully not a reference
	}
}

func (mt MintTransaction) AsRow() []string {
	return []string{
		mt.Date.Format(DateOutputFormat),
		mt.Description,
		fmt.Sprint(mt.Amount),
		mt.TransactionType,
		mt.Category,
		mt.Account,
		mt.Notes,
	}
}

func (mt MintTransaction) Display() string {
	return fmt.Sprintf(
		"{ Date: '%s', Description: '%s', Amount: %s, TransactionType: '%s', Category: '%s', Account: '%s', Notes: '%s' }",
		mt.Date.Format(DateOutputFormat),
		mt.Description,
		mt.Amount,
		mt.TransactionType,
		mt.Category,
		mt.Account,
		mt.Notes,
	)
}
