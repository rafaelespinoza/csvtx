package csvtx

import "testing"

func TestYnabAsRow(t *testing.T) {
	tables := []struct {
		yt       YnabTransaction
		expected []string
	}{
		{
			YnabTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "foo", Amount(-4321)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "43.21", ""},
		},
		{
			YnabTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "foo", Amount(4321)},
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

func TestXflow(t *testing.T) {
	tables := []struct {
		amount    Amount
		direction string
		expected  string
	}{
		{Amount(-100), "in", ""},
		{Amount(-100), "out", "1.00"},
		{Amount(100), "in", "1.00"},
		{Amount(100), "out", ""},
	}

	for _, test := range tables {
		actual := xflow(test.amount, test.direction)

		if actual != test.expected {
			t.Errorf(
				"\ninputs: (%d, %s)\nactual: %s\nexpected: %s\n",
				test.amount,
				test.direction,
				actual,
				test.expected,
			)
		}

	}
}
