package csvtx

import (
	"fmt"
	"strings"
)

func (t Transaction) String() string {
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
