package csvtx

import (
	"fmt"
	"os"
	"testing"
)

func TestWriteAcctFiles(t *testing.T) {
	input := "./fixtures/mint.csv"

	ReadParseMint(input, func(m []MintTransaction) {
		expectedOutputs := []string{"checking.csv", "personal-savings.csv", "credit.csv"}
		WriteAcctFiles(m)

		for _, f := range expectedOutputs {
			if _, err := os.Stat(f); err != nil {
				t.Errorf("expected file %s to be created\n", f)
			} else {
				os.Remove(f)
				fmt.Printf("removed %s\n", f)
			}
		}
	})
}
