package ynab

import (
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

var Headers = []string{
	"Date", "Payee", "Category", "Memo", "Outflow", "Inflow",
}

type Transaction struct {
	Date     time.Time
	Payee    string
	Category string
	Memo     string
	Amount   entity.AmountSubunits
}

func (t Transaction) AsRow() []string {
	return []string{
		t.Date.Format(entity.DateOutputFormat),
		t.Payee,
		t.Category,
		t.Memo,
		formatAmount(t.Amount, false),
		formatAmount(t.Amount, true),
	}
}

func formatAmount(amt entity.AmountSubunits, inflow bool) string {
	negative := amt < 0

	if (negative && !inflow) || (!negative && inflow) {
		return amt.String()
	}

	return ""
}
