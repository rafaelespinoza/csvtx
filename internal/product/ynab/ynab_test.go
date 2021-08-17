package ynab

import (
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func TestAsRow(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)
	tables := []struct {
		yt       Transaction
		expected []string
	}{
		{
			Transaction{date, "Joe's Diner", "Restaurants", "foo", entity.AmountSubunits(-4321)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "43.21", ""},
		},
		{
			Transaction{date, "Joe's Diner", "Restaurants", "foo", entity.AmountSubunits(4321)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "", "43.21"},
		},
	}

	for _, test := range tables {
		output := test.yt.AsRow()
		var expected string

		for i, actual := range output {
			expected = test.expected[i]

			if actual != expected {
				t.Errorf("\nactual: %s\nexpected: %s\n", actual, expected)
			}
		}

	}
}

func TestFormatAmount(t *testing.T) {
	tables := []struct {
		amount   entity.AmountSubunits
		inflow   bool
		expected string
	}{
		{entity.AmountSubunits(-100), true, ""},
		{entity.AmountSubunits(-100), false, "1.00"},
		{entity.AmountSubunits(100), true, "1.00"},
		{entity.AmountSubunits(100), false, ""},
	}

	for _, test := range tables {
		actual := formatAmount(test.amount, test.inflow)

		if actual != test.expected {
			t.Errorf(
				"\ninputs: (%d, %t)\nactual: %s\nexpected: %s\n",
				test.amount,
				test.inflow,
				actual,
				test.expected,
			)
		}
	}
}
