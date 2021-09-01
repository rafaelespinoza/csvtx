package convert

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestVenmoToYNAB(t *testing.T) {
	baseDir, err := os.MkdirTemp("", strings.Replace(t.Name(), "/", "_", -1)+"_*")
	if err != nil {
		t.Fatal(err)
	}

	params := Params{
		Infile:  "testdata/venmo.csv",
		Outdir:  baseDir,
		LogDest: &testLogger{t},
	}
	err = VenmoToYNAB(params)
	if err != nil {
		t.Fatal(err)
	}

	got, err := readAllOutput(path.Join(baseDir, "venmo.csv"))
	if err != nil {
		t.Fatal(err)
	}
	expectedData := [][]string{
		{"Date", "Payee", "Category", "Memo", "Outflow", "Inflow"},
		{"06/12/2021", "Chuck", "", "Roam 🍔", "18.00", ""},
		{"06/16/2021", "Biff", "", "dog stuff", "", "138.00"},
		{"06/19/2021", "Lorax", "", "Dinner", "", "30.00"},
		{"06/23/2021", "Chuck", "", "Pizza", "37.00", ""},
		{"07/30/2021", "Biff", "", "class", "60.00", ""},
		{"08/19/2021", "Mork", "", "Stuff for camping trip Costco", "132.00", ""},
		{"08/26/2021", "Garth", "", "⛽ 🏕", "40.00", ""},
		{"08/29/2021", "Mork", "", "Pretzels", "", "50.00"},
		{"08/30/2021", "Garth", "", "rounding test", "8.78", ""},
	}
	if len(got) != len(expectedData) {
		t.Fatalf("wrong number of output rows; got %d, expected %d", len(got), len(expectedData))
	}

	for i, row := range got {
		testOutputRow(t, row, expectedData[i])
	}
}
