package csvtx

import (
	"testing"
	"time"
)

func TestReadParseMint(t *testing.T) {
	filename := "./fixtures/mint.csv"
	ReadParseMint(filename, func(m []MintTransaction) {
		if len(m) < 1 {
			t.Errorf("\nShould be able to read file & pass some populated data to this callback. Instead got empty results\n")
		}
	})
}

func TestParseLine(t *testing.T) {
	tables := []struct {
		input    []string
		expected MintTransaction
	}{
		{
			[]string{"04/01/2018", "Joe's Diner", "JOES DINER SF CA", "43.21", "debit", "Restaurants", "Checking", "labels", "some notes"},
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(-4321), "debit"},
		},
		{
			[]string{"04/01/2018", "Fancy Clothes Inc", "Conglomerate", "250.00", "debit", "Shopping", "Credit", "", "pants"},
			MintTransaction{th_makeDate(), "Fancy Clothes Inc", "Shopping", "Credit", "pants", Amount(-25000), "debit"},
		},
		{
			[]string{"04/01/2018", "Paycheck", "ACME", "1000000000.00", "credit", "Income", "Checking", "label", "monthly income"},
			MintTransaction{th_makeDate(), "Paycheck", "Income", "Checking", "monthly income", Amount(100000000000), "credit"},
		},
	}

	for _, test := range tables {
		actual := parseLine(test.input)

		if !actual.Equal(test.expected) {
			t.Errorf("\nactual:   %v\nexpected: %v\n", actual, test.expected)
		}
	}
}

func TestParseMoney(t *testing.T) {
	tables := []struct {
		cell       string
		isNegative bool
		expected   Amount
	}{
		{"", false, Amount(0)},
		{"", true, Amount(0)},
		{"0", false, Amount(0)},
		{"0", true, Amount(0)},
		{"0.01", false, Amount(1)},
		{"0.01", true, Amount(-1)},
		{"0.99", false, Amount(99)},
		{"0.99", true, Amount(-99)},
		{"12", false, Amount(1200)},
		{"12", true, Amount(-1200)},
		{"12.34", false, Amount(1234)},
		{"12.34", true, Amount(-1234)},
		{"567", false, Amount(56700)},
		{"567", true, Amount(-56700)},
	}

	for _, test := range tables {
		actual := parseMoney(test.cell, test.isNegative)

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
			t.Errorf("%s != %s\n", actual, test.expected)
		}
	}
}
