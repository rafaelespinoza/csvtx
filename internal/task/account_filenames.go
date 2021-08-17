package task

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/rafaelespinoza/csvtx/internal/product/mint"
)

// accountFilenames associates an account name (ie: checking, savings, credit)
// to a filename that has transaction data for that account.
type accountFilenames map[string]string

func newAccountFilenames(mints *[]mint.Transaction) accountFilenames {
	uniqAcctNames := accountFilenames{}

	var txAcctType string

	for _, tx := range *mints {
		txAcctType = tx.Account

		if _, ok := uniqAcctNames[txAcctType]; !ok {
			uniqAcctNames[txAcctType] = accountToFilename(txAcctType)
		}
	}

	return uniqAcctNames
}

var (
	space     = regexp.MustCompile(`\W`)
	uppercase = regexp.MustCompile(`[A-Z]`)
)

func accountToFilename(acctName string) string {
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
