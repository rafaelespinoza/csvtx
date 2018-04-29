package csvtx

import (
	"time"
)

type YnabTransaction struct {
	Date     time.Time
	CheckNum int
	Payee    string
	Memo     string
	Amount   Amount // cents
}
