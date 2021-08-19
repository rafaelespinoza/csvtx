package cmd

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/rafaelespinoza/alf"
	"github.com/rafaelespinoza/csvtx/internal/task"
)

func makeConvert(cmdName string) alf.Directive {
	var params task.Params

	out := alf.Delegator{
		Description: "convert csv to csv data for another service/program",
		Flags:       flag.NewFlagSet(cmdName, flag.ExitOnError),
	}
	out.Flags.StringVar(&params.Infile, "i", "", "path to input csv")
	out.Flags.StringVar(&params.Outdir, "o", "/tmp", "path to output directory")

	out.Subs = map[string]alf.Directive{
		"mint-to-ynab": &alf.Command{
			Description: "convert Mint CSV data to YNAB4 data",
			Setup: func(p flag.FlagSet) *flag.FlagSet {
				p.Usage = func() {
					fmt.Fprintf(p.Output(), `Usage: %s %s

Description:

Flags:

`, _Bin, cmdName)
					p.PrintDefaults()
				}
				return &p
			},
			Run: func(ctx context.Context) error {
				return task.MintToYNAB(params)
			},
		},
	}

	out.Flags.Usage = func() {
		fmt.Fprintf(out.Flags.Output(), `Description: convert CSV to CSV

Subcommands:

	%v

Flags:

`, strings.Join(out.DescribeSubcommands(), "\n\t"))
		out.Flags.PrintDefaults()
	}

	return &out
}
