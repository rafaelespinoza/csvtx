package task

import (
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
	"github.com/rafaelespinoza/csvtx/internal/product/mint"
)

const pathToFixtures = "../../fixtures"

func TestReadParseMint(t *testing.T) {
	ReadParseMint(pathToFixtures+"/mint.csv", func(m []mint.Transaction) {
		if len(m) < 1 {
			t.Errorf("Should be able to read file & pass some populated data to this callback. Instead got empty results")
		}
	})
}

func TestParseLine(t *testing.T) {
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
		actual := parseLine(test.input)

		if !actual.Equal(test.expected) {
			t.Error("not equal")
			t.Logf("actual:   %v", actual)
			t.Logf("expected: %v", test.expected)
		}
	}
}

func TestParseMoney(t *testing.T) {
	tables := []struct {
		cell     string
		negative bool
		expected entity.AmountSubunits
	}{
		{"", false, entity.AmountSubunits(0)},
		{"", true, entity.AmountSubunits(0)},
		{"0", false, entity.AmountSubunits(0)},
		{"0", true, entity.AmountSubunits(0)},
		{"0.01", false, entity.AmountSubunits(1)},
		{"0.01", true, entity.AmountSubunits(-1)},
		{"0.99", false, entity.AmountSubunits(99)},
		{"0.99", true, entity.AmountSubunits(-99)},
		{"12", false, entity.AmountSubunits(1200)},
		{"12", true, entity.AmountSubunits(-1200)},
		{"12.34", false, entity.AmountSubunits(1234)},
		{"12.34", true, entity.AmountSubunits(-1234)},
		{"567", false, entity.AmountSubunits(56700)},
		{"567", true, entity.AmountSubunits(-56700)},
	}

	for _, test := range tables {
		actual := parseMoney(test.cell, test.negative)

		if actual != test.expected {
			t.Errorf("%v != %v\n", actual, test.expected)
		}
	}
}

func TestParseDate(t *testing.T) {
	tables := []struct {
		input    string
		expected time.Time
	}{
		{"04/01/2018", time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)},
		{"4/01/2018", time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, test := range tables {
		actual := parseDate(test.input)

		if !actual.Equal(test.expected) {
			t.Errorf("%v != %v\n", actual, test.expected)
		}
	}
}

func TestFallbackStr(t *testing.T) {
	tables := []struct {
		input    string
		fallback string
		expected string
	}{
		{"", "default val", "default val"},
		{"foo", "default val", "foo"},
	}

	for _, test := range tables {
		actual := fallbackStr(test.input, test.fallback)

		if actual != test.expected {
			t.Errorf("wrong value; got %q, expected %q", actual, test.expected)
		}
	}
}
