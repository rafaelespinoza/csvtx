package csvtx

import "testing"

func TestNewAccountNameMap(t *testing.T) {
	tables := []struct {
		input    *[]MintTransaction
		expected AccountNameMap
	}{
		{
			&[]MintTransaction{
				MintTransaction{th_makeDate(), "payee", "cat", "Savings", "n", Amount(51), "credit"},
				MintTransaction{th_makeDate(), "payee", "cat", "Checking", "n", Amount(4321), "debit"},
				MintTransaction{th_makeDate(), "payee", "cat", "Credit", "n", Amount(25000), "debit"},
				MintTransaction{th_makeDate(), "payee", "cat", "Credit", "n", Amount(100000000000), "credit"},
				MintTransaction{th_makeDate(), "payee", "cat", "Credit", "n", Amount(1000000000), "debit"},
				MintTransaction{th_makeDate(), "payee", "cat", "Checking", "n", Amount(100000000000), "credit"},
			},
			AccountNameMap{
				"Checking": "checking.csv",
				"Credit":   "credit.csv",
				"Savings":  "savings.csv",
			},
		},
	}

	for _, test := range tables {
		actual := newAccountNameMap(test.input)

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

func TestAcctToFileName(t *testing.T) {
	tables := []struct {
		input    string
		expected string
	}{
		{"FOO BAR", "foo-bar.csv"},
		{"Personal savings", "personal-savings.csv"},
	}

	for _, test := range tables {
		actual := acctToFileName(test.input)

		if actual != test.expected {
			t.Errorf("%s != %s\n", actual, test.expected)
		}
	}
}
