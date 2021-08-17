package task

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
	"github.com/rafaelespinoza/csvtx/internal/product/mint"
	"github.com/rafaelespinoza/csvtx/internal/product/ynab"
)

func TestWriteAcctFiles(t *testing.T) {
	ReadParseMint(pathToFixtures+"/mint.csv", func(m []mint.Transaction) {
		expectedOutputs := []string{"checking.csv", "personal-savings.csv", "credit.csv"}
		WriteAcctFiles(m)

		for _, f := range expectedOutputs {
			if _, err := os.Stat(f); err != nil {
				t.Errorf("expected file %s to be created\n", f)
			} else {
				os.Remove(f)
				fmt.Printf("removed %s\n", f)
			}
		}
	})
}

func TestMintToYnab(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)

	tables := []struct {
		mt       mint.Transaction
		expected ynab.Transaction
	}{
		{
			mint.Transaction{Date: date, Description: "Joe's Diner", Category: "Restaurants", Account: "Checking", Notes: "some notes", Amount: entity.AmountSubunits(4321), TransactionType: "debit"},
			ynab.Transaction{Date: date, Payee: "Joe's Diner", Category: "Restaurants", Memo: "some notes", Amount: entity.AmountSubunits(4321)},
		},
		{
			mint.Transaction{Date: date, Description: "Fancy Clothes Inc", Category: "Shopping", Account: "Credit", Notes: "pants", Amount: entity.AmountSubunits(25000), TransactionType: "debit"},
			ynab.Transaction{Date: date, Payee: "Fancy Clothes Inc", Category: "Shopping", Memo: "pants", Amount: entity.AmountSubunits(25000)},
		},
		{
			mint.Transaction{Date: date, Description: "Paycheck", Category: "Income", Account: "Checking", Notes: "monthly income", Amount: entity.AmountSubunits(100000000000), TransactionType: "credit"},
			ynab.Transaction{Date: date, Payee: "Paycheck", Category: "Income", Memo: "monthly income", Amount: entity.AmountSubunits(100000000000)},
		},
	}

	for _, test := range tables {
		output := mintToYnab(test.mt)
		expected := test.expected

		if output != expected {
			t.Errorf("\nactual:   %v\nexpected: %v\n", output, expected)
		}
	}
}
