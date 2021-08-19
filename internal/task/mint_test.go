package task

import (
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func TestParseMintRow(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)

	tables := []struct {
		input    []string
		expected entity.Mint
	}{
		{
			[]string{"04/01/2018", "Joe's Diner", "JOES DINER SF CA", "43.21", "debit", "Restaurants", "Checking", "labels", "some notes"},
			entity.Mint{Date: date, Description: "Joe's Diner", Category: "Restaurants", Account: "Checking", Notes: "some notes", Amount: entity.AmountSubunits(-4321), TransactionType: "debit"},
		},
		{
			[]string{"04/01/2018", "Fancy Clothes Inc", "Conglomerate", "250.00", "debit", "Shopping", "Credit", "", "pants"},
			entity.Mint{Date: date, Description: "Fancy Clothes Inc", Category: "Shopping", Account: "Credit", Notes: "pants", Amount: entity.AmountSubunits(-25000), TransactionType: "debit"},
		},
		{
			[]string{"04/01/2018", "Paycheck", "ACME", "1000000000.00", "credit", "Income", "Checking", "label", "monthly income"},
			entity.Mint{Date: date, Description: "Paycheck", Category: "Income", Account: "Checking", Notes: "monthly income", Amount: entity.AmountSubunits(100000000000), TransactionType: "credit"},
		},
	}

	for i, test := range tables {
		actual, err := parseMintRow(test.input)
		if err != nil {
			t.Fatalf("test %d; %v", i, err)
		}

		if actual.Amount != test.expected.Amount {
			t.Errorf("test %d; wrong Amount; got %d, expected %d", i, actual.Amount, test.expected.Amount)
		}
		if !actual.Date.Equal(test.expected.Date) {
			t.Errorf("test %d; wrong Date; got %s, expected %s", i, actual.Date, test.expected.Date)
		}
		if actual.Description != test.expected.Description {
			t.Errorf("test %d; wrong Description; got %s, expected %s", i, actual.Description, test.expected.Description)
		}
		if actual.Category != test.expected.Category {
			t.Errorf("test %d; wrong Category; got %s, expected %s", i, actual.Category, test.expected.Category)
		}
		if actual.Account != test.expected.Account {
			t.Errorf("test %d; wrong Account; got %s, expected %s", i, actual.Account, test.expected.Account)
		}
		if actual.Notes != test.expected.Notes {
			t.Errorf("test %d; wrong Notes; got %s, expected %s", i, actual.Notes, test.expected.Notes)
		}
		if actual.TransactionType != test.expected.TransactionType {
			t.Errorf("test %d; wrong TransactionType; got %s, expected %s", i, actual.TransactionType, test.expected.TransactionType)
		}
	}
}
