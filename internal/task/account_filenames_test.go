package task

import (
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
	"github.com/rafaelespinoza/csvtx/internal/product/mint"
)

func TestNewAccountFilenames(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)

	tables := []struct {
		input    *[]mint.Transaction
		expected accountFilenames
	}{
		{
			&[]mint.Transaction{
				{Date: date, Description: "payee", Category: "cat", Account: "Savings", Notes: "n", Amount: entity.AmountSubunits(51), TransactionType: "credit"},
				{Date: date, Description: "payee", Category: "cat", Account: "Checking", Notes: "n", Amount: entity.AmountSubunits(4321), TransactionType: "debit"},
				{Date: date, Description: "payee", Category: "cat", Account: "Credit", Notes: "n", Amount: entity.AmountSubunits(25000), TransactionType: "debit"},
				{Date: date, Description: "payee", Category: "cat", Account: "Credit", Notes: "n", Amount: entity.AmountSubunits(100000000000), TransactionType: "credit"},
				{Date: date, Description: "payee", Category: "cat", Account: "Credit", Notes: "n", Amount: entity.AmountSubunits(1000000000), TransactionType: "debit"},
				{Date: date, Description: "payee", Category: "cat", Account: "Checking", Notes: "n", Amount: entity.AmountSubunits(100000000000), TransactionType: "credit"},
			},
			accountFilenames{
				"Checking": "checking.csv",
				"Credit":   "credit.csv",
				"Savings":  "savings.csv",
			},
		},
	}

	for _, test := range tables {
		actual := newAccountFilenames(test.input)

		n := len(test.expected)

		if len(actual) != n {
			t.Errorf("output does not have same num of key/val pairs as expected (%d != %d)\n", len(actual), n)
		}

		for k, v := range actual {
			if v != test.expected[k] {
				t.Errorf("values do not match up (%s != %s)\n", v, test.expected[k])
			}
		}
	}
}

func TestAccountToFilename(t *testing.T) {
	tables := []struct {
		input    string
		expected string
	}{
		{"FOO BAR", "foo-bar.csv"},
		{"Personal savings", "personal-savings.csv"},
	}

	for _, test := range tables {
		actual := accountToFilename(test.input)

		if actual != test.expected {
			t.Errorf("%s != %s\n", actual, test.expected)
		}
	}
}
