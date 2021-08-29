package convert

import (
	"encoding/csv"
	"os"
	"testing"
	"time"

	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func TestParseMoney(t *testing.T) {
	tables := []struct {
		cell     string
		negative bool
		expected entity.AmountSubunits
	}{
		{"", false, entity.AmountSubunits(0)},
		{"", true, entity.AmountSubunits(0)},
		{"0", false, entity.AmountSubunits(0)},
		{"0", true, entity.AmountSubunits(0)},
		{"0.01", false, entity.AmountSubunits(1)},
		{"0.01", true, entity.AmountSubunits(-1)},
		{"0.99", false, entity.AmountSubunits(99)},
		{"0.99", true, entity.AmountSubunits(-99)},
		{"12", false, entity.AmountSubunits(1200)},
		{"12", true, entity.AmountSubunits(-1200)},
		{"12.34", false, entity.AmountSubunits(1234)},
		{"12.34", true, entity.AmountSubunits(-1234)},
		{"567", false, entity.AmountSubunits(56700)},
		{"567", true, entity.AmountSubunits(-56700)},
	}

	for _, test := range tables {
		actual, err := parseMoney(test.cell, test.negative)
		if err != nil {
			t.Fatal(err)
		}

		if actual != test.expected {
			t.Errorf("%v != %v\n", actual, test.expected)
		}
	}
}

func TestParseDate(t *testing.T) {
	tables := []struct {
		input    string
		expected time.Time
	}{
		{"04/01/2018", time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)},
		{"4/01/2018", time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, test := range tables {
		actual, err := parseDate(test.input)
		if err != nil {
			t.Fatal(err)
		}

		if !actual.Equal(test.expected) {
			t.Errorf("%v != %v\n", actual, test.expected)
		}
	}
}

func TestYNAB(t *testing.T) {
	date := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)
	runTest := func(t *testing.T, in entity.YNAB, expected []string) {
		t.Helper()

		output := ynabAsRow(in)
		if len(output) != len(expected) {
			t.Fatalf("wrong number of values; got %d, expected %d", len(output), len(expected))
		}

		for i, got := range output {
			if got != expected[i] {
				t.Errorf("wrong value at [%d]; got %q, expected %q", i, got, expected[i])
			}
		}
	}

	t.Run("outflow", func(t *testing.T) {
		runTest(
			t,
			entity.YNAB{Date: date, Payee: "Joe's Diner", Category: "Restaurants", Memo: "foo", Amount: entity.AmountSubunits(-4321)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "43.21", ""},
		)

		runTest(
			t,
			entity.YNAB{Date: date, Payee: "Joe's Diner", Category: "Restaurants", Memo: "foo", Amount: entity.AmountSubunits(-100)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "1.00", ""},
		)
	})

	t.Run("inflow", func(t *testing.T) {
		runTest(
			t,
			entity.YNAB{Date: date, Payee: "Joe's Diner", Category: "Restaurants", Memo: "foo", Amount: entity.AmountSubunits(4321)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "", "43.21"},
		)

		runTest(
			t,
			entity.YNAB{Date: date, Payee: "Joe's Diner", Category: "Restaurants", Memo: "foo", Amount: entity.AmountSubunits(100)},
			[]string{"04/01/2018", "Joe's Diner", "Restaurants", "foo", "", "1.00"},
		)
	})
}

type testLogger struct{ t *testing.T }

func (w *testLogger) Write(in []byte) (n int, e error) {
	w.t.Logf("%s", in)
	n = len(in)
	return
}

func readAllOutput(filename string) (out [][]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	out, err = reader.ReadAll()
	return
}

func testOutputRow(t *testing.T, actual, expected []string) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("wrong number of data values; got %d, expected %d", len(actual), len(expected))
	}

	for j, val := range actual {
		if val != expected[j] {
			t.Errorf("wrong value at column[%d]; got %q, expected %q", j, val, expected[j])
		}
	}
}
