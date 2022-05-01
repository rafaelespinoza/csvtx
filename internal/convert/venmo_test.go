package convert

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
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

func TestReadParseVenmo(t *testing.T) {
	const okInputHeader = `Account Statement - (@FooBar) - June 1st to September 1st 2021 ,,,,,,,,,,,,,,,,,,
Account Activity,,,,,,,,,,,,,,,,,,
,ID,Datetime,Type,Status,Note,From,To,Amount (total),Amount (tip),Amount (fee),Funding Source,Destination,Beginning Balance,Ending Balance,Statement Period Venmo Fees,Terminal Location,Year to Date Venmo Fees,Disclaimer
`

	t.Run("ok", func(t *testing.T) {
		const (
			inputZeroItems = okInputHeader + `,,,,,,,,,,,,,$100.00,,,,,
`
			inputMultipleItems = okInputHeader + `,,,,,,,,,,,,,$100.00,,,,,
,1111111111111111111,2021-06-12T03:13:29,Payment,Complete,Roam 🍔,Lorax,Chuck,- $18.00,,,Venmo balance,,,,,Venmo,,
,4444444444444444444,2021-06-23T20:06:31,Charge,Complete,Pizza,Chuck,Lorax,- $37.00,,,Venmo balance,,,,,Venmo,,
`
			// Can handle another file, possibly with its own headers.
			inputMultipleFiles = inputMultipleItems + okInputHeader + `,5555555555555555555,2021-07-30T23:18:13,Payment,Complete,class,Lorax,Biff,- $60.00,,,Venmo balance,,,,,Venmo,,
`
		)

		tests := []struct {
			name     string
			input    string
			sortAsc  bool
			expected []*entity.Venmo
		}{
			{name: "sort asc - no items", input: inputZeroItems, sortAsc: true, expected: []*entity.Venmo{}},
			{
				name:    "sort asc - multiple items",
				input:   inputMultipleItems,
				sortAsc: true,
				expected: []*entity.Venmo{
					{
						Datetime:        time.Date(2021, time.June, 12, 03, 13, 29, 00, time.UTC),
						TransactionType: "Payment",
						Note:            "Roam 🍔",
						From:            "Lorax",
						To:              "Chuck",
						Amount:          -1800,
					},
					{
						Datetime:        time.Date(2021, time.June, 23, 20, 06, 31, 00, time.UTC),
						TransactionType: "Charge",
						Note:            "Pizza",
						From:            "Chuck",
						To:              "Lorax",
						Amount:          -3700,
					},
				},
			},
			{
				name:    "sort asc - multiple files",
				input:   inputMultipleFiles,
				sortAsc: true,
				expected: []*entity.Venmo{
					{
						Datetime:        time.Date(2021, time.June, 12, 03, 13, 29, 00, time.UTC),
						TransactionType: "Payment",
						Note:            "Roam 🍔",
						From:            "Lorax",
						To:              "Chuck",
						Amount:          -1800,
					},
					{
						Datetime:        time.Date(2021, time.June, 23, 20, 06, 31, 00, time.UTC),
						TransactionType: "Charge",
						Note:            "Pizza",
						From:            "Chuck",
						To:              "Lorax",
						Amount:          -3700,
					},
					{
						Datetime:        time.Date(2021, time.July, 30, 23, 18, 13, 00, time.UTC),
						TransactionType: "Payment",
						Note:            "class",
						From:            "Lorax",
						To:              "Biff",
						Amount:          -6000,
					},
				},
			},
			{name: "sort desc - no items", input: inputZeroItems, sortAsc: false, expected: []*entity.Venmo{}},
			{
				name:    "sort desc - multiple items",
				input:   inputMultipleItems,
				sortAsc: false,
				expected: []*entity.Venmo{
					{
						Datetime:        time.Date(2021, time.June, 23, 20, 06, 31, 00, time.UTC),
						TransactionType: "Charge",
						Note:            "Pizza",
						From:            "Chuck",
						To:              "Lorax",
						Amount:          -3700,
					},
					{
						Datetime:        time.Date(2021, time.June, 12, 03, 13, 29, 00, time.UTC),
						TransactionType: "Payment",
						Note:            "Roam 🍔",
						From:            "Lorax",
						To:              "Chuck",
						Amount:          -1800,
					},
				},
			},
			{
				name:    "sort desc - multiple files",
				input:   inputMultipleFiles,
				sortAsc: false,
				expected: []*entity.Venmo{
					{
						Datetime:        time.Date(2021, time.July, 30, 23, 18, 13, 00, time.UTC),
						TransactionType: "Payment",
						Note:            "class",
						From:            "Lorax",
						To:              "Biff",
						Amount:          -6000,
					},
					{
						Datetime:        time.Date(2021, time.June, 23, 20, 06, 31, 00, time.UTC),
						TransactionType: "Charge",
						Note:            "Pizza",
						From:            "Chuck",
						To:              "Lorax",
						Amount:          -3700,
					},
					{
						Datetime:        time.Date(2021, time.June, 12, 03, 13, 29, 00, time.UTC),
						TransactionType: "Payment",
						Note:            "Roam 🍔",
						From:            "Lorax",
						To:              "Chuck",
						Amount:          -1800,
					},
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				actual, err := ReadParseVenmo(strings.NewReader(test.input), test.sortAsc)
				if err != nil {
					t.Fatal(err)
				}

				if len(actual) != len(test.expected) {
					t.Fatalf("wrong number of items; got %d, expected %d", len(actual), len(test.expected))
				}

				for i, got := range actual {
					exp := test.expected[i]

					if !got.Datetime.Equal(exp.Datetime) {
						t.Errorf("item %d, wrong Datetime; got %s, expected %s", i, got.Datetime.Format(time.RFC3339), exp.Datetime.Format(time.RFC3339))
					}

					if got.TransactionType != exp.TransactionType {
						t.Errorf("item %d, wrong TransactionType; got %s, expected %s", i, got.TransactionType, exp.TransactionType)
					}

					if got.Note != exp.Note {
						t.Errorf("item %d, wrong Note; got %s, expected %s", i, got.Note, exp.Note)
					}

					if got.From != exp.From {
						t.Errorf("item %d, wrong From; got %s, expected %s", i, got.From, exp.From)
					}

					if got.To != exp.To {
						t.Errorf("item %d, wrong To; got %s, expected %s", i, got.To, exp.To)
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
				input:             okInputHeader + `,1111111111111111111,bad_date,Payment,Complete,Roam 🍔,Lorax,Chuck,- $18.00,,,Venmo balance,,,,,Venmo,,`,
				expErrMsgContains: "bad_date",
			},
			{
				name:              "bad Amount",
				input:             okInputHeader + `,1111111111111111111,2021-06-12T03:13:29,Payment,Complete,Roam 🍔,Lorax,Chuck,- bad_amount,,,Venmo balance,,,,,Venmo,,`,
				expErrMsgContains: "bad_amount",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := ReadParseVenmo(strings.NewReader(test.input), true)
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
