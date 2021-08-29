package convert

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestMechanicsBankToYNAB(t *testing.T) {
	runTest := func(t *testing.T, filename string) {
		baseDir, err := os.MkdirTemp("", strings.Replace(t.Name(), "/", "_", -1)+"_*")
		if err != nil {
			t.Fatal(err)
		}

		params := Params{
			Infile:  filename,
			Outdir:  baseDir,
			LogDest: &testLogger{t},
		}
		err = MechanicsBankToYNAB(params)
		if err != nil {
			t.Fatal(err)
		}

		got, err := readAllOutput(path.Join(baseDir, "mechanicsbank.csv"))
		if err != nil {
			t.Fatal(err)
		}
		expectedData := [][]string{
			{"Date", "Payee", "Category", "Memo", "Outflow", "Inflow"},
			{"08/13/2021", "71214 FOOBAR DIR DEP", "", "", "", "1234.56"},
			{"07/28/2021", "POWER COMPANY WEB ONLINE", "", "", "40.47", ""},
			{"07/03/2021", "STOP ITEM CHARGE(S)", "", "", "25.00", ""},
			{"06/21/2021", "RENT", "", "", "1111.00", ""},
		}
		if len(got) != len(expectedData) {
			t.Fatalf("wrong number of output rows; got %d, expected %d", len(got), len(expectedData))
		}

		for i, row := range got {
			testOutputRow(t, row, expectedData[i])
		}
	}

	t.Run("includes balance", func(t *testing.T) {
		runTest(t, "testdata/mechanicsbank.balance.csv")
	})
	t.Run("excludes balance", func(t *testing.T) {
		runTest(t, "testdata/mechanicsbank.csv")
	})
}
