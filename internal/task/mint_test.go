package task

import (
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
	"github.com/rafaelespinoza/csvtx/internal/product/mint"
)

func TestParseMintRow(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)

	tables := []struct {
		input    []string
		expected mint.Transaction
	}{
		{
			[]string{"04/01/2018", "Joe's Diner", "JOES DINER SF CA", "43.21", "debit", "Restaurants", "Checking", "labels", "some notes"},
			mint.Transaction{Date: date, Description: "Joe's Diner", Category: "Restaurants", Account: "Checking", Notes: "some notes", Amount: entity.AmountSubunits(-4321), TransactionType: "debit"},
		},
		{
			[]string{"04/01/2018", "Fancy Clothes Inc", "Conglomerate", "250.00", "debit", "Shopping", "Credit", "", "pants"},
			mint.Transaction{Date: date, Description: "Fancy Clothes Inc", Category: "Shopping", Account: "Credit", Notes: "pants", Amount: entity.AmountSubunits(-25000), TransactionType: "debit"},
		},
		{
			[]string{"04/01/2018", "Paycheck", "ACME", "1000000000.00", "credit", "Income", "Checking", "label", "monthly income"},
			mint.Transaction{Date: date, Description: "Paycheck", Category: "Income", Account: "Checking", Notes: "monthly income", Amount: entity.AmountSubunits(100000000000), TransactionType: "credit"},
		},
	}

	for _, test := range tables {
		actual, err := parseMintRow(test.input)
		if err != nil {
			t.Fatal(err)
		}

		if !actual.Equal(test.expected) {
			t.Error("not equal")
			t.Logf("actual:   %v", actual)
			t.Logf("expected: %v", test.expected)
		}
	}
}
