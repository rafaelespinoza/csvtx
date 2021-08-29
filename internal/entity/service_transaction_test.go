package entity_test

import (
	"testing"

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
