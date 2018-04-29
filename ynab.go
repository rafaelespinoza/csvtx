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
		yt.Date.Format(DateOutputFormat),
		fmt.Sprint(yt.CheckNum),
		yt.Payee,
		yt.Memo,
		fmt.Sprint(yt.Amount),
	}
}
func (yt YnabTransaction) Display() string {
	return fmt.Sprintf(
		"{ Date: '%s', CheckNum: %d, Payee: '%s', Memo: '%s', Amount: '%s' }",
		yt.Date.Format(DateOutputFormat),
		yt.CheckNum,
		yt.Payee,
		yt.Memo,
		yt.Amount,
	)
}
