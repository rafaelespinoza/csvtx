package entity_test

import (
	"testing"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func TestAmount(t *testing.T) {
	tables := []struct {
		amount   entity.AmountSubunits
		expected string
	}{
		{entity.AmountSubunits(-4321), "43.21"},
		{entity.AmountSubunits(4321), "43.21"},
		{entity.AmountSubunits(-99), "0.99"},
		{entity.AmountSubunits(99), "0.99"},
		{entity.AmountSubunits(0), "0.00"},
	}

	for _, test := range tables {
		actual := test.amount.String()

		if actual != test.expected {
			t.Errorf("\nactual: %s\nexpected: %s\n", actual, test.expected)
		}
	}
}
