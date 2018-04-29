package csvtx

import (
	"bytes"
	"regexp"
	"strings"
)

var space = regexp.MustCompile(`\W`)
var uppercase = regexp.MustCompile(`[A-Z]`)

type AccountNameMap map[string]string

func NewAccountNameMap(mints *[]MintTransaction) AccountNameMap {
	uniqAcctNames := AccountNameMap{}

	var txAcctType string

	for _, tx := range *mints {
		txAcctType = tx.Account

		if _, ok := uniqAcctNames[txAcctType]; !ok {
			uniqAcctNames[txAcctType] = acctToFileName(txAcctType)
		}
	}

	return uniqAcctNames
}

func acctToFileName(acctName string) string {
	var buf bytes.Buffer
	var ch string

	n := len(acctName)

	for i := 0; i < n; i++ {
		ch = string(acctName[i])

		if space.MatchString(ch) {
			buf.WriteString("-")
		} else if uppercase.MatchString(ch) {
			lower := strings.ToLower(ch)
			buf.WriteString(lower)
		} else {
			buf.WriteString(ch)
		}
	}

	return buf.String() + ".csv"
}
