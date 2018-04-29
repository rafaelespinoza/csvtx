package csvtx

import (
	"fmt"
	"time"
)

var YnabHeader = []string{
	"Date", "Payee", "Category", "Memo", "Outflow", "Inflow",
}

type YnabTransaction struct {
	Date     time.Time
	Payee    string
	Category string
	Memo     string
	Amount   Amount // cents
}

func (yt YnabTransaction) AsRow() []string {
	return []string{
		yt.Date.Format(DateOutputFormat),
		yt.Payee,
		yt.Category,
		yt.Memo,
		xflow(yt.Amount, "out"),
		xflow(yt.Amount, "in"),
	}
}

func xflow(amt Amount, direction string) string {
	isNegative := amt.isNegative()
	isInflow := direction == "in"

	if (isNegative && isInflow) || (!isNegative && !isInflow) {
		return ""
	} else if (isNegative && !isInflow) || (!isNegative && isInflow) {
		return amt.String()
	} else {
		return fmt.Sprintf("y'all fucked up: %d, %s", amt, direction)
	}
}

func (yt YnabTransaction) Display() string {
	amt := yt.Amount.String()

	if yt.Amount.isNegative() {
		amt = "(" + amt + ")"
	}

	return fmt.Sprintf(
		"{ Date: '%s', Payee: '%s', Category: %s, Memo: '%s', Amount: '%s' }",
		yt.Date.Format(DateOutputFormat),
		yt.Payee,
		yt.Category,
		yt.Memo,
		amt,
	)
}
