package csvtx

import (
	"fmt"
	"time"
)

type MintTransaction struct {
	Date            time.Time
	Description     string
	Category        string
	Account         string
	Notes           string
	Amount          Amount // cents
	transactionType string // [debit, credit]

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
		mt.Category,
		mt.Account,
		mt.Notes,
	}
}

func (mt MintTransaction) isNegative() (bool, error) {
	tt := mt.transactionType

	if tt == TransactionTypes[0] {
		return true, nil
	} else if tt == TransactionTypes[1] {
		return false, nil
	} else {
		err := fmt.Errorf("%T: %v has invalid transaction type : %s", mt, mt, tt)
		return false, err
	}
}

func (mt MintTransaction) Display() string {
	return fmt.Sprintf(
		"{ Date: '%s', Description: '%s', Amount: %s, Category: '%s', Account: '%s', Notes: '%s' }",
		mt.Date.Format(DateOutputFormat),
		mt.Description,
		mt.Amount,
		mt.Category,
		mt.Account,
		mt.Notes,
	)
}
