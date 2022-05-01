package convert

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
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
			{"06/11/2021", "ROUND TEST", "", "", "71.82", ""},
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

func TestReadParseMechanicsBank(t *testing.T) {
	const okInputHeader = `AccountName : Checking
"Account Number : 1234"
"Date Range : 07/27/2021-08/16/2021"
Transaction Number,Date,Description,Memo,Amount Debit,Amount Credit,Check Number,Fees
`

	t.Run("ok", func(t *testing.T) {
		const (
			inputZeroItems     = okInputHeader + ``
			inputMultipleItems = okInputHeader + `"20210813000000[-8:PST]*1234.56*501**71214 FOO BAR DIR DEP",08/13/2021,"71214 FOOBAR DIR DEP","",,1234.56,,0.00
"20210280000000[-8:PST]*-40.47*0**POWER COMPANY WEB ONLINE",07/28/2021,"POWER COMPANY WEB ONLINE","",-40.47,,,0.00
`
			// Can handle another file, possibly with its own headers.
			inputMultipleFiles = inputMultipleItems + okInputHeader + `"20210621000000[-8:PST]*-1111.00*0**RENT",06/21/2021,"RENT","",-1111.00,,,0.00
`
		)

		tests := []struct {
			name     string
			input    string
			sortAsc  bool
			expected []*entity.MechanicsBank
		}{
			{name: "sort asc - no items", input: inputZeroItems, sortAsc: true, expected: []*entity.MechanicsBank{}},
			{
				name:    "sort asc - multiple items",
				input:   inputMultipleItems,
				sortAsc: true,
				expected: []*entity.MechanicsBank{
					{Date: time.Date(2021, time.July, 28, 0, 0, 0, 0, time.UTC), AmountDebit: 4047, Description: "POWER COMPANY WEB ONLINE"},
					{Date: time.Date(2021, time.August, 13, 0, 0, 0, 0, time.UTC), AmountCredit: 123456, Description: "71214 FOOBAR DIR DEP"},
				},
			},
			{
				name:    "sort asc - multiple files",
				input:   inputMultipleFiles,
				sortAsc: true,
				expected: []*entity.MechanicsBank{
					{Date: time.Date(2021, time.June, 21, 0, 0, 0, 0, time.UTC), AmountDebit: 111100, Description: "RENT"},
					{Date: time.Date(2021, time.July, 28, 0, 0, 0, 0, time.UTC), AmountDebit: 4047, Description: "POWER COMPANY WEB ONLINE"},
					{Date: time.Date(2021, time.August, 13, 0, 0, 0, 0, time.UTC), AmountCredit: 123456, Description: "71214 FOOBAR DIR DEP"},
				},
			},
			{name: "sort desc - no items", input: inputZeroItems, sortAsc: false, expected: []*entity.MechanicsBank{}},
			{
				name:    "sort desc - multiple items",
				input:   inputMultipleItems,
				sortAsc: false,
				expected: []*entity.MechanicsBank{
					{Date: time.Date(2021, time.August, 13, 0, 0, 0, 0, time.UTC), AmountCredit: 123456, Description: "71214 FOOBAR DIR DEP"},
					{Date: time.Date(2021, time.July, 28, 0, 0, 0, 0, time.UTC), AmountDebit: 4047, Description: "POWER COMPANY WEB ONLINE"},
				},
			},
			{
				name:    "sort desc - multiple files",
				input:   inputMultipleFiles,
				sortAsc: false,
				expected: []*entity.MechanicsBank{
					{Date: time.Date(2021, time.August, 13, 0, 0, 0, 0, time.UTC), AmountCredit: 123456, Description: "71214 FOOBAR DIR DEP"},
					{Date: time.Date(2021, time.July, 28, 0, 0, 0, 0, time.UTC), AmountDebit: 4047, Description: "POWER COMPANY WEB ONLINE"},
					{Date: time.Date(2021, time.June, 21, 0, 0, 0, 0, time.UTC), AmountDebit: 111100, Description: "RENT"},
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				actual, err := ReadParseMechanicsBank(strings.NewReader(test.input), test.sortAsc)
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

					if got.AmountDebit != exp.AmountDebit {
						t.Errorf("item %d, wrong AmountDebit; got %d, expected %d", i, got.AmountDebit, exp.AmountDebit)
					}

					if got.AmountCredit != exp.AmountCredit {
						t.Errorf("item %d, wrong AmountCredit; got %d, expected %d", i, got.AmountCredit, exp.AmountCredit)
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
				input:             okInputHeader + `"20210813000000[-8:PST]*1234.56*501**71214 FOO BAR DIR DEP",bad_date,"71214 FOOBAR DIR DEP","",,1234.56,,0.00`,
				expErrMsgContains: "bad_date",
			},
			{
				name:              "bad AmountDebit",
				input:             okInputHeader + `"20210813000000[-8:PST]*1234.56*501**71214 FOO BAR DIR DEP",08/13/2021,"71214 FOOBAR DIR DEP","",bad_debit_amount,,,0.00`,
				expErrMsgContains: "bad_debit_amount",
			},
			{
				name:              "bad AmountCredit",
				input:             okInputHeader + `"20210813000000[-8:PST]*1234.56*501**71214 FOO BAR DIR DEP",08/13/2021,"71214 FOOBAR DIR DEP","",,bad_credit_amount,,0.00`,
				expErrMsgContains: "bad_credit_amount",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := ReadParseMechanicsBank(strings.NewReader(test.input), true)
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
