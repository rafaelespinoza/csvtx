package csvtx

import "testing"

func TestMintAsYnabTx(t *testing.T) {
	tables := []struct {
		mt       MintTransaction
		expected YnabTransaction
	}{
		{
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "debit"},
			YnabTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "some notes", Amount(4321)},
		},
		{
			MintTransaction{th_makeDate(), "Fancy Clothes Inc", "Shopping", "Credit", "pants", Amount(25000), "debit"},
			YnabTransaction{th_makeDate(), "Fancy Clothes Inc", "Shopping", "pants", Amount(25000)},
		},
		{
			MintTransaction{th_makeDate(), "Paycheck", "Income", "Checking", "monthly income", Amount(100000000000), "credit"},
			YnabTransaction{th_makeDate(), "Paycheck", "Income", "monthly income", Amount(100000000000)},
		},
	}

	for _, test := range tables {
		output := test.mt.asYnabTx()
		expected := test.expected

		if output != expected {
			t.Errorf("\nactual:   %v\nexpected: %v\n", output, expected)
		}
	}

}

func TestMintTransactionAsRow(t *testing.T) {
	tables := []struct {
		mt       MintTransaction
		expected []string
	}{
		{
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "debit"},
			[]string{"04/01/2018", "Joe's Diner", "43.21", "Restaurants", "Checking", "some notes"},
		},
		{
			MintTransaction{th_makeDate(), "Fancy Clothes Inc", "Shopping", "Credit", "pants", Amount(25000), "debit"},
			[]string{"04/01/2018", "Fancy Clothes Inc", "250.00", "Shopping", "Credit", "pants"},
		},
		{
			MintTransaction{th_makeDate(), "Paycheck", "Income", "Checking", "monthly income", Amount(100000000000), "credit"},
			[]string{"04/01/2018", "Paycheck", "1000000000.00", "Income", "Checking", "monthly income"},
		},
	}

	for _, test := range tables {
		output := test.mt.AsRow()
		var expected string

		for i, actual := range output {
			expected = test.expected[i]

			if actual != expected {
				t.Errorf("\nactual:   %s\nexpected: %s\n", actual, expected)
			}
		}

	}
}

func TestMintEqual(t *testing.T) {
	tables := []struct {
		mt       MintTransaction
		mu       MintTransaction
		expected bool
	}{
		{
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "debit"},
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "debit"},
			true,
		},
		{
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "debit"},
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(-4321), "debit"},
			false,
		},
		{
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "debit"},
			MintTransaction{th_makeDate().AddDate(0, 0, 1), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "debit"},
			false,
		},
		{
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "debit"},
			MintTransaction{th_makeDate(), "Joe's Diner", "Restaurants", "Checking", "some notes", Amount(4321), "credit"},
			false,
		},
	}

	for _, test := range tables {
		actual := test.mt.Equal(test.mu)
		if actual != test.expected {
			t.Errorf("expected %v to be equal to %v\n", test.mt, test.mu)
		}
	}
}
