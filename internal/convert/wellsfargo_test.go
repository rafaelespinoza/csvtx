package convert

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
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

func TestReadParseWellsFargo(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		const (
			inputZeroItems = ``

			inputMultipleItems = `"08/16/2021","692.72","*","","ONLINE ACH PAYMENT - THANK YOU"
"08/09/2021","-22.31","*","","COFFEE SHOP WWW.EXAMPLE.COM"
`
			// Can handle another file, possibly with its own headers.
			inputMultipleFiles = inputMultipleItems + `"06/21/2021","-64.10","*","","float round test"
`
		)

		tests := []struct {
			name     string
			input    string
			sortAsc  bool
			expected []*entity.WellsFargo
		}{
			{name: "sort asc - no items", input: inputZeroItems, sortAsc: true, expected: []*entity.WellsFargo{}},
			{
				name:    "sort asc - multiple items",
				input:   inputMultipleItems,
				sortAsc: true,
				expected: []*entity.WellsFargo{
					{Date: time.Date(2021, time.August, 9, 0, 0, 0, 0, time.UTC), Amount: -2231, Description: "COFFEE SHOP WWW.EXAMPLE.COM"},
					{Date: time.Date(2021, time.August, 16, 0, 0, 0, 0, time.UTC), Amount: 69272, Description: "ONLINE ACH PAYMENT - THANK YOU"},
				},
			},
			{
				name:    "sort asc - multiple inputMultipleFiles",
				input:   inputMultipleFiles,
				sortAsc: true,
				expected: []*entity.WellsFargo{
					{Date: time.Date(2021, time.June, 21, 0, 0, 0, 0, time.UTC), Amount: -6410, Description: "float round test"},
					{Date: time.Date(2021, time.August, 9, 0, 0, 0, 0, time.UTC), Amount: -2231, Description: "COFFEE SHOP WWW.EXAMPLE.COM"},
					{Date: time.Date(2021, time.August, 16, 0, 0, 0, 0, time.UTC), Amount: 69272, Description: "ONLINE ACH PAYMENT - THANK YOU"},
				},
			},
			{name: "sort desc - no items", input: inputZeroItems, sortAsc: false, expected: []*entity.WellsFargo{}},
			{
				name:    "sort desc - multiple items",
				input:   inputMultipleItems,
				sortAsc: false,
				expected: []*entity.WellsFargo{
					{Date: time.Date(2021, time.August, 16, 0, 0, 0, 0, time.UTC), Amount: 69272, Description: "ONLINE ACH PAYMENT - THANK YOU"},
					{Date: time.Date(2021, time.August, 9, 0, 0, 0, 0, time.UTC), Amount: -2231, Description: "COFFEE SHOP WWW.EXAMPLE.COM"},
				},
			},
			{
				name:    "sort desc - multiple files",
				input:   inputMultipleFiles,
				sortAsc: false,
				expected: []*entity.WellsFargo{
					{Date: time.Date(2021, time.August, 16, 0, 0, 0, 0, time.UTC), Amount: 69272, Description: "ONLINE ACH PAYMENT - THANK YOU"},
					{Date: time.Date(2021, time.August, 9, 0, 0, 0, 0, time.UTC), Amount: -2231, Description: "COFFEE SHOP WWW.EXAMPLE.COM"},
					{Date: time.Date(2021, time.June, 21, 0, 0, 0, 0, time.UTC), Amount: -6410, Description: "float round test"},
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				actual, err := ReadParseWellsFargo(strings.NewReader(test.input), test.sortAsc)
				if err != nil {
					t.Fatal(err)
				}

				if len(actual) != len(test.expected) {
					t.Fatalf("wrong number of items; got %d, expected %d", len(actual), len(test.expected))
				}

				for i, got := range actual {
					exp := test.expected[i]

					if !got.Date.Equal(exp.Date) {
						t.Errorf("item %d, wrong Date; got %s, expected %s", i, got.Date.Format(time.RFC3339), exp.Date.Format(time.RFC3339))
					}

					if got.Description != exp.Description {
						t.Errorf("item %d, wrong Description; got %s, expected %s", i, got.Description, exp.Description)
					}

					if got.Amount != exp.Amount {
						t.Errorf("item %d, wrong Amount; got %d, expected %d", i, got.Amount, exp.Amount)
					}
				}
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		tests := []struct {
			name              string
			input             string
			expErrMsgContains string
		}{
			{
				name:              "bad Date",
				input:             `"bad_date","692.72","*","","ONLINE ACH PAYMENT - THANK YOU"`,
				expErrMsgContains: "bad_date",
			},
			{
				name:              "bad Amount",
				input:             `"08/16/2021","bad_amount","*","","ONLINE ACH PAYMENT - THANK YOU"`,
				expErrMsgContains: "bad_amount",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := ReadParseWellsFargo(strings.NewReader(test.input), true)
				if err == nil {
					t.Fatal("expected an error, got nil")
				}

				got := err.Error()
				if !strings.Contains(got, test.expErrMsgContains) {
					t.Errorf("expected error message (%q) to contain %q", got, test.expErrMsgContains)
				}
			})
		}
	})
}
