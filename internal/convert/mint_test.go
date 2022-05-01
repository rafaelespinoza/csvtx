package convert

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
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
			{"01/18/2018", "Off by 1", "test", "round test", "73.21", ""},
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
			{"01/17/2018", "Off by 1", "test", "round test", "", "39.55"},
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

func TestReadParseMint(t *testing.T) {
	const okInputHeader = `"Date","Description","Original Description","Amount","Transaction Type","Category","Account Name","Labels","Notes"
`

	t.Run("ok", func(t *testing.T) {
		const (
			inputZeroItems     = okInputHeader + ``
			inputMultipleItems = okInputHeader + `"4/25/2018","Transfer Checking","Online Banking Transfer","1098.76","credit","Transfer","PERSONAL SAVINGS","",""
"2/20/2018","Transfer","","1234.56","debit","Transfer","Checking","",""
`
			// Can handle another file, possibly with its own headers.
			inputMultipleFiles = inputMultipleItems + okInputHeader + `"01/20/2018","Fancy Clothes Inc","Conglomerate","250.00","debit","Shopping","Credit","","pants"
`
		)

		tests := []struct {
			name     string
			input    string
			sortAsc  bool
			expected []*entity.Mint
		}{
			{name: "sort asc - no items", input: inputZeroItems, sortAsc: true, expected: []*entity.Mint{}},
			{
				name:    "sort asc - multiple items",
				input:   inputMultipleItems,
				sortAsc: true,
				expected: []*entity.Mint{
					{
						Date:            time.Date(2018, time.February, 20, 0, 0, 0, 0, time.UTC),
						Description:     "Transfer",
						TransactionType: "debit",
						Category:        "Transfer",
						Account:         "Checking",
						Notes:           "",
						Amount:          -123456,
					},
					{
						Date:            time.Date(2018, time.April, 25, 0, 0, 0, 0, time.UTC),
						Description:     "Transfer Checking",
						TransactionType: "credit",
						Category:        "Transfer",
						Account:         "PERSONAL SAVINGS",
						Notes:           "",
						Amount:          109876,
					},
				},
			},
			{
				name:    "sort asc - multiple files",
				input:   inputMultipleFiles,
				sortAsc: true,
				expected: []*entity.Mint{
					{
						Date:            time.Date(2018, time.January, 20, 0, 0, 0, 0, time.UTC),
						Description:     "Fancy Clothes Inc",
						TransactionType: "debit",
						Category:        "Shopping",
						Account:         "Credit",
						Notes:           "pants",
						Amount:          -25000,
					},
					{
						Date:            time.Date(2018, time.February, 20, 0, 0, 0, 0, time.UTC),
						Description:     "Transfer",
						TransactionType: "debit",
						Category:        "Transfer",
						Account:         "Checking",
						Notes:           "",
						Amount:          -123456,
					},
					{
						Date:            time.Date(2018, time.April, 25, 0, 0, 0, 0, time.UTC),
						Description:     "Transfer Checking",
						TransactionType: "credit",
						Category:        "Transfer",
						Account:         "PERSONAL SAVINGS",
						Notes:           "",
						Amount:          109876,
					},
				},
			},
			{name: "sort desc - no items", input: inputZeroItems, sortAsc: false, expected: []*entity.Mint{}},
			{
				name:    "sort desc - multiple items",
				input:   inputMultipleItems,
				sortAsc: false,
				expected: []*entity.Mint{
					{
						Date:            time.Date(2018, time.April, 25, 0, 0, 0, 0, time.UTC),
						Description:     "Transfer Checking",
						TransactionType: "credit",
						Category:        "Transfer",
						Account:         "PERSONAL SAVINGS",
						Notes:           "",
						Amount:          109876,
					},
					{
						Date:            time.Date(2018, time.February, 20, 0, 0, 0, 0, time.UTC),
						Description:     "Transfer",
						TransactionType: "debit",
						Category:        "Transfer",
						Account:         "Checking",
						Notes:           "",
						Amount:          -123456,
					},
				},
			},
			{
				name:    "sort desc - multiple files",
				input:   inputMultipleFiles,
				sortAsc: false,
				expected: []*entity.Mint{
					{
						Date:            time.Date(2018, time.April, 25, 0, 0, 0, 0, time.UTC),
						Description:     "Transfer Checking",
						TransactionType: "credit",
						Category:        "Transfer",
						Account:         "PERSONAL SAVINGS",
						Notes:           "",
						Amount:          109876,
					},
					{
						Date:            time.Date(2018, time.February, 20, 0, 0, 0, 0, time.UTC),
						Description:     "Transfer",
						TransactionType: "debit",
						Category:        "Transfer",
						Account:         "Checking",
						Notes:           "",
						Amount:          -123456,
					},
					{
						Date:            time.Date(2018, time.January, 20, 0, 0, 0, 0, time.UTC),
						Description:     "Fancy Clothes Inc",
						TransactionType: "debit",
						Category:        "Shopping",
						Account:         "Credit",
						Notes:           "pants",
						Amount:          -25000,
					},
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				actual, err := ReadParseMint(strings.NewReader(test.input), test.sortAsc)
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

					if got.TransactionType != exp.TransactionType {
						t.Errorf("item %d, wrong TransactionType; got %s, expected %s", i, got.TransactionType, exp.TransactionType)
					}

					if got.Category != exp.Category {
						t.Errorf("item %d, wrong Category; got %s, expected %s", i, got.Category, exp.Category)
					}

					if got.Account != exp.Account {
						t.Errorf("item %d, wrong Account; got %s, expected %s", i, got.Account, exp.Account)
					}

					if got.Notes != exp.Notes {
						t.Errorf("item %d, wrong Notes; got %s, expected %s", i, got.Notes, exp.Notes)
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
				input:             okInputHeader + `"bad_date","Joe's Diner","JOES DINER SF CA","43.21","debit","Restaurants","Checking","labels","some notes"`,
				expErrMsgContains: "bad_date",
			},
			{
				name:              "bad TransactionType",
				input:             okInputHeader + `"1/20/2018","Paycheck","ACME","12.34","bad_transaction_type","Income","Checking","label","monthly income"`,
				expErrMsgContains: "bad_transaction_type",
			},
			{
				name:              "bad Amount",
				input:             okInputHeader + `"1/20/2018","Paycheck","ACME","bad_amount","credit","Income","Checking","label","monthly income"`,
				expErrMsgContains: "bad_amount",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := ReadParseMint(strings.NewReader(test.input), true)
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
