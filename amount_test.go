package csvtx

import "testing"

func TestAmountString(t *testing.T) {
	tables := []struct {
		amount   Amount
		expected string
	}{
		{Amount(-4321), "43.21"},
		{Amount(4321), "43.21"},
		{Amount(-99), "0.99"},
		{Amount(99), "0.99"},
		{Amount(0), "0.00"},
	}

	for _, test := range tables {
		actual := test.amount.String()

		if actual != test.expected {
			t.Errorf("\nactual: %s\nexpected: %s\n", actual, test.expected)
		}
	}
}

func TestIsNegative(t *testing.T) {
	tables := []struct {
		amount   Amount
		expected bool
	}{
		{Amount(-100), true},
		{Amount(100), false},
		{Amount(0), false},
	}

	for _, test := range tables {
		actual := test.amount.isNegative()

		if actual != test.expected {
			t.Errorf(
				"\nactual %t, expected %t",
				actual,
				test.expected,
			)
		}

	}
}
