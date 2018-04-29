package csvtx

import (
	"fmt"
	"time"
)

type YnabTransaction struct {
	Date     time.Time
	CheckNum int
	Payee    string
	Memo     string
	Amount   Amount // cents
}

func (yt YnabTransaction) AsRow() []string {
	return []string{
		yt.Date.Format("2006-01-02"),
		fmt.Sprint(yt.CheckNum),
		yt.Payee,
		yt.Memo,
		fmt.Sprint(yt.Amount),
	}
}
