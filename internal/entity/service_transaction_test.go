package entity_test

import (
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func TestMint(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		got, err := entity.Mint{TransactionType: "debit"}.Negative()
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Errorf("wrong Negative; got %t, expected %t", got, true)
		}
	})

	t.Run("false", func(t *testing.T) {
		got, err := entity.Mint{TransactionType: "credit"}.Negative()
		if err != nil {
			t.Fatal(err)
		}
		if got {
			t.Errorf("wrong Negative; got %t, expected %t", got, false)
		}
	})

	t.Run("error", func(t *testing.T) {
		_, err := entity.Mint{TransactionType: "invalid"}.Negative()
		if err == nil {
			t.Fatal("should reject invalid inputs")
		}
	})
}

func TestYNAB(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)
	runTest := func(t *testing.T, in entity.YNAB, expected []string) {
		t.Helper()

		output := in.AsRow()
		if len(output) != len(expected) {
			t.Fatalf("wrong number of values; got %d, expected %d", len(output), len(expected))
		}

		for i, got := range output {
			if got != expected[i] {
				t.Errorf("wrong value at [%d]; got %q, expected %q", i, got, expected[i])
			}
		}
	}

	t.Run("outflow", func(t *testing.T) {
		runTest(
			t,
			entity.YNAB{date, "Joe's Diner", "Restaurants", "foo", entity.AmountSubunits(-4321)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "43.21", ""},
		)

		runTest(
			t,
			entity.YNAB{date, "Joe's Diner", "Restaurants", "foo", entity.AmountSubunits(-100)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "1.00", ""},
		)
	})

	t.Run("inflow", func(t *testing.T) {
		runTest(
			t,
			entity.YNAB{date, "Joe's Diner", "Restaurants", "foo", entity.AmountSubunits(4321)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "", "43.21"},
		)

		runTest(
			t,
			entity.YNAB{date, "Joe's Diner", "Restaurants", "foo", entity.AmountSubunits(100)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "", "1.00"},
		)
	})
}
