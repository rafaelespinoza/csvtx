package task

import (
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

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
		actual, err := parseMoney(test.cell, test.negative)
		if err != nil {
			t.Fatal(err)
		}

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
		actual, err := parseDate(test.input)
		if err != nil {
			t.Fatal(err)
		}

		if !actual.Equal(test.expected) {
			t.Errorf("%v != %v\n", actual, test.expected)
		}
	}
}
