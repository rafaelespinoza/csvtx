package cmd

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/rafaelespinoza/alf"
	"github.com/rafaelespinoza/csvtx/internal/convert"
)

func makeConvert(cmdName string) alf.Directive {
	var params struct {
		Infile string
		Outdir string
		From   string
	}

	out := alf.Command{
		Description: "convert csv to csv data for YNAB4 import",
		Setup: func(p flag.FlagSet) *flag.FlagSet {
			flags := flag.NewFlagSet(cmdName, flag.ExitOnError)
			flags.StringVar(&params.Infile, "i", "", "path to input csv")
			flags.StringVar(&params.Outdir, "o", "/tmp", "path to output directory")
			flags.StringVar(
				&params.From,
				"from",
				"",
				fmt.Sprintf("product/service of input file; one of %v", fromServices),
			)

			flags.Usage = func() {
				fmt.Fprintf(flags.Output(), `Usage: %s %s -from from -i infile -o outdir

Description:

	Converts an input CSV from a source service to a CSV format ready for import
	into YNAB 4 classic edition. The source service must be one of:

		%v

	If flag i is empty, then it reads from standard input.

Flags:

`, _Bin, cmdName, strings.Join(fromServices, "\n\t\t"))
				flags.PrintDefaults()
			}

			return flags
		},
		Run: func(_ context.Context) error {
			p := convert.Params{
				Infile: params.Infile,
				Outdir: params.Outdir,
			}
			switch params.From {
			case serviceMechanicsBank:
				return convert.MechanicsBankToYNAB(p)
			case serviceMint:
				return convert.MintToYNAB(p)
			case serviceVenmo:
				return convert.VenmoToYNAB(p)
			case serviceWellsFargo:
				return convert.WellsFargoToYNAB(p)
			default:
				return fmt.Errorf("unknown source service %q", params.From)
			}
		},
	}

	return &out
}
