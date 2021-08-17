package mint

import (
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func TestAsRow(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)
	tables := []struct {
		mt       Transaction
		expected []string
	}{
		{
			Transaction{date, "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(4321), "debit"},
			[]string{"04/01/2018", "Joe's Diner", "43.21", "Restaurants", "Checking", "some notes"},
		},
		{
			Transaction{date, "Fancy Clothes Inc", "Shopping", "Credit", "pants", entity.AmountSubunits(25000), "debit"},
			[]string{"04/01/2018", "Fancy Clothes Inc", "250.00", "Shopping", "Credit", "pants"},
		},
		{
			Transaction{date, "Paycheck", "Income", "Checking", "monthly income", entity.AmountSubunits(100000000000), "credit"},
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

func TestEqual(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)
	tables := []struct {
		mt       Transaction
		mu       Transaction
		expected bool
	}{
		{
			Transaction{date, "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(4321), "debit"},
			Transaction{date, "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(4321), "debit"},
			true,
		},
		{
			Transaction{date, "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(4321), "debit"},
			Transaction{date, "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(-4321), "debit"},
			false,
		},
		{
			Transaction{date, "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(4321), "debit"},
			Transaction{date.AddDate(0, 0, 1), "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(4321), "debit"},
			false,
		},
		{
			Transaction{date, "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(4321), "debit"},
			Transaction{date, "Joe's Diner", "Restaurants", "Checking", "some notes", entity.AmountSubunits(4321), "credit"},
			false,
		},
	}

	for _, test := range tables {
		actual := test.mt.Equal(test.mu)
		if actual != test.expected {
			t.Errorf("expected %v to be equal to %v", test.mt, test.mu)
		}
	}
}

func TestNegative(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		got, err := Transaction{TransactionType: "debit"}.Negative()
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Errorf("wrong Negative; got %t, expected %t", got, true)
		}
	})

	t.Run("false", func(t *testing.T) {
		got, err := Transaction{TransactionType: "credit"}.Negative()
		if err != nil {
			t.Fatal(err)
		}
		if got {
			t.Errorf("wrong Negative; got %t, expected %t", got, false)
		}
	})

	t.Run("error", func(t *testing.T) {
		_, err := Transaction{TransactionType: "invalid"}.Negative()
		if err == nil {
			t.Fatal("should reject invalid inputs")
		}
	})
}
