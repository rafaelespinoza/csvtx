package mint

import (
	"fmt"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

type Transaction struct {
	Date            time.Time
	Description     string
	Category        string
	Account         string
	Notes           string
	Amount          entity.AmountSubunits
	TransactionType string // [debit, credit]

	// ignoring these columns from mint csv file:
	// "Original Description", "Labels"
}

func (t *Transaction) Equal(other Transaction) bool {
	if t.Amount != other.Amount {
		return false
	}

	if !t.Date.Equal(other.Date) {
		return false
	}

	theseAttrs := t.strAttrs()
	otherAttrs := other.strAttrs()

	for i, a := range theseAttrs {
		if a != otherAttrs[i] {
			return false
		}
	}

	return true
}

func (t *Transaction) strAttrs() []string {
	return []string{
		t.Description, t.Category, t.Account, t.Notes, t.TransactionType,
	}
}

func (t Transaction) AsRow() []string {
	return []string{
		t.Date.Format(entity.DateOutputFormat),
		t.Description,
		fmt.Sprint(t.Amount),
		t.Category,
		t.Account,
		t.Notes,
	}
}

func (t Transaction) Negative() (bool, error) {
	switch tt := t.TransactionType; tt {
	case "debit":
		return true, nil
	case "credit":
		return false, nil
	default:
		return false, fmt.Errorf("invalid transaction type %q", tt)
	}
}
