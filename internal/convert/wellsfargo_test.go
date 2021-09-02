package convert

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestWellsFargo(t *testing.T) {
	baseDir, err := os.MkdirTemp("", strings.Replace(t.Name(), "/", "_", -1)+"_*")
	if err != nil {
		t.Fatal(err)
	}

	params := Params{
		Infile:  "testdata/wellsfargo.csv",
		Outdir:  baseDir,
		LogDest: &testLogger{t},
	}
	err = WellsFargoToYNAB(params)
	if err != nil {
		t.Fatal(err)
	}

	got, err := readAllOutput(path.Join(baseDir, "wellsfargo.csv"))
	if err != nil {
		t.Fatal(err)
	}
	expectedData := [][]string{
		{"Date", "Payee", "Category", "Memo", "Outflow", "Inflow"},
		{"08/16/2021", "ONLINE ACH PAYMENT - THANK YOU", "", "", "", "692.72"},
		{"08/09/2021", "COFFEE SHOP WWW.EXAMPLE.COM", "", "", "22.31", ""},
		{"08/02/2021", "GROCERY STORE SPRINGFIELD USA", "", "", "66.43", ""},
		{"07/26/2021", "BRIDGE TOLL", "", "", "5.00", ""},
		{"07/19/2021", "KFC", "", "", "213.62", ""},
		{"07/12/2021", "SUBSCRIPTION", "", "", "9.99", ""},
		{"07/05/2021", "PET FOOD", "", "", "16.30", ""},
		{"06/28/2021", "web services www.example.coWA", "", "", "0.78", ""},
		{"06/21/2021", "float round test", "", "", "64.10", ""},
	}
	if len(got) != len(expectedData) {
		t.Fatalf("wrong number of output rows; got %d, expected %d", len(got), len(expectedData))
	}

	for i, row := range got {
		testOutputRow(t, row, expectedData[i])
	}
}
