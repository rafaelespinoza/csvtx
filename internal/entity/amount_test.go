package entity

import "testing"

func TestAmount(t *testing.T) {
	tables := []struct {
		amount   AmountSubunits
		expected string
	}{
		{AmountSubunits(-4321), "43.21"},
		{AmountSubunits(4321), "43.21"},
		{AmountSubunits(-99), "0.99"},
		{AmountSubunits(99), "0.99"},
		{AmountSubunits(0), "0.00"},
	}

	for _, test := range tables {
		actual := test.amount.String()

		if actual != test.expected {
			t.Errorf("\nactual: %s\nexpected: %s\n", actual, test.expected)
		}
	}
}
