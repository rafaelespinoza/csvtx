package cmd

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rafaelespinoza/alf"
	"github.com/rafaelespinoza/csvtx/internal/convert"
	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func makeNormalizeCSV(parentName, name string) alf.Directive {
	var params *normalizeParams

	out := alf.Command{
		Description: "read CSV from stdin and print normalized data to stdout as CSV",
		Setup: func(p flag.FlagSet) (flags *flag.FlagSet) {
			fullname := parentName + " " + name
			flags, params = initNormalizeSubcmd(fullname)

			flags.Usage = func() {
				fmt.Fprintf(flags.Output(), `Usage: %s [flags] < input.csv

Description:

	Pipe in csv data, parse it, print the parsed version to stdout.
	Here, "normalize" means that dates and amounts have consistent presentation.

	The from argument must be one of:

		%v

Flags:

`, fullname, strings.Join(fromServices, "\n\t\t"))
				flags.PrintDefaults()
			}

			return flags
		},
		Run: func(_ context.Context) (err error) {
			in, out := os.Stdin, os.Stdout

			switch params.from {
			case serviceMechanicsBank:
				var data []*entity.MechanicsBank
				data, err = convert.ReadParseMechanicsBank(in, params.asc)
				if err != nil {
					return
				}

				records := make([][]string, len(data))
				for i, dat := range data {
					records[i] = mechanicsBankToCSV(dat, params.datelayout)
				}
				err = writeAllCSV(records, out)
			case serviceMint:
				var data []*entity.Mint
				data, err = convert.ReadParseMint(in, params.asc)
				if err != nil {
					return
				}

				records := make([][]string, len(data))
				for i, dat := range data {
					records[i] = mintToCSV(dat, params.datelayout)
				}
				err = writeAllCSV(records, out)
			case serviceVenmo:
				var data []*entity.Venmo
				data, err = convert.ReadParseVenmo(in, params.asc)
				if err != nil {
					return
				}

				records := make([][]string, len(data))
				for i, dat := range data {
					records[i] = venmoToCSV(dat, params.datelayout)
				}
				err = writeAllCSV(records, out)
			case serviceWellsFargo:
				var data []*entity.WellsFargo
				data, err = convert.ReadParseWellsFargo(in, params.asc)
				if err != nil {
					return
				}

				records := make([][]string, len(data))
				for i, dat := range data {
					records[i] = wellsFargoToCSV(dat, params.datelayout)
				}
				err = writeAllCSV(records, out)
			default:
				err = fmt.Errorf("unknown source service %q", params.from)
			}

			return
		},
	}

	return &out
}

func writeAllCSV(records [][]string, out io.Writer) (err error) {
	err = csv.NewWriter(out).WriteAll(records)
	return
}

func mechanicsBankToCSV(in *entity.MechanicsBank, dateLayout string) []string {
	return []string{
		in.Date.Format(dateLayout),
		in.Description,
		in.Memo,
		in.AmountDebit.String(),
		in.AmountCredit.String(),
		in.Balance.String(),
		fmt.Sprintf("%d", in.CheckNumber),
		in.Fees.String(),
	}
}

func mintToCSV(in *entity.Mint, dateLayout string) []string {
	return []string{
		in.Date.Format(dateLayout),
		in.Description,
		in.Category,
		in.Account,
		in.Notes,
		in.Amount.String(),
		in.TransactionType,
	}
}

func venmoToCSV(in *entity.Venmo, dateLayout string) []string {
	return []string{
		in.Datetime.Format(dateLayout),
		in.TransactionType,
		in.Note,
		in.From,
		in.To,
		in.Amount.String(),
	}
}

func wellsFargoToCSV(in *entity.WellsFargo, dateLayout string) []string {
	return []string{
		in.Date.Format(dateLayout),
		in.Amount.String(),
		in.Description,
	}
}
