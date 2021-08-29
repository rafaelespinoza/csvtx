package convert

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestMintToYNAB(t *testing.T) {
	baseDir, err := os.MkdirTemp("", strings.Replace(t.Name(), "/", "_", -1)+"_*")
	if err != nil {
		t.Fatal(err)
	}

	params := Params{
		Infile:  "testdata/mint.csv",
		Outdir:  baseDir,
		LogDest: &testLogger{t},
	}
	err = MintToYNAB(params)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string][][]string{
		"checking.csv": {
			{"Date", "Payee", "Category", "Memo", "Outflow", "Inflow"},
			{"04/25/2018", "Transfer Savings", "Transfer", "", "1098.76", ""},
			{"03/15/2018", "Transfer Savings", "Transfer", "", "789.10", ""},
			{"02/20/2018", "Transfer", "Transfer", "", "1234.56", ""},
			{"02/05/2018", "Payment", "Transfer", "", "1234.56", ""},
			{"02/05/2018", "Joe's Diner", "Restaurants", "some notes", "43.21", ""},
			{"01/20/2018", "Paycheck", "Income", "monthly income", "", "1000000000.00"},
		},
		"personal-savings.csv": {
			{"Date", "Payee", "Category", "Memo", "Outflow", "Inflow"},
			{"04/25/2018", "Transfer Checking", "Transfer", "", "", "1098.76"},
			{"04/10/2018", "Interest Paid", "Interest Income", "", "", "0.06"},
			{"03/15/2018", "Transfer to Checking", "Transfer", "", "789.10", ""},
			{"02/20/2018", "Deposit", "Transfer", "", "", "1234.56"},
		},
		"credit.csv": {
			{"Date", "Payee", "Category", "Memo", "Outflow", "Inflow"},
			{"01/20/2018", "Fancy Clothes Inc", "Shopping", "pants", "250.00", ""},
		},
	}
	for filename, expectedData := range tests {
		filename = path.Join(baseDir, filename)
		got, err := readAllOutput(filename)
		if err != nil {
			t.Fatalf("could not read or parse file %q; %v", filename, err)
		}
		if len(got) != len(expectedData) {
			t.Fatalf("wrong number of output rows; got %d, expected %d", len(got), len(expectedData))
		}

		for i, row := range got {
			testOutputRow(t, row, expectedData[i])
		}
	}
}
