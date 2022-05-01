package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/rafaelespinoza/alf"
)

const defaultDateLayout = "2006-01-02"

func makeNormalize(cmdName string) alf.Directive {
	fullname := _Bin + " " + cmdName

	out := alf.Delegator{
		Description: "read csv from stdin and print to stdout",
		Flags:       flag.NewFlagSet(fullname, flag.ExitOnError),
		Subs: map[string]alf.Directive{
			"csv":  makeNormalizeCSV(fullname, "csv"),
			"json": makeNormalizeJSON(fullname, "json"),
		},
	}

	out.Flags.Usage = func() {
		fmt.Fprintf(out.Flags.Output(), `Description:

	Pipe in csv data, parse it, print the parsed version to stdout.
	Here, "normalize" means that dates and amounts have consistent presentation.

	Each output format has its own subcommand.

Subcommands:

	These may have their own set of flags. Put them after the subcommand.

	%v

Flags:

`, strings.Join(out.DescribeSubcommands(), "\n\t"))
		out.Flags.PrintDefaults()
	}

	return &out
}

type normalizeParams struct {
	datelayout string
	from       string
	asc        bool
}

func initNormalizeSubcmd(fullname string) (flags *flag.FlagSet, params *normalizeParams) {
	flags = flag.NewFlagSet(fullname, flag.ExitOnError)
	params = &normalizeParams{}

	flags.StringVar(
		&params.datelayout,
		"datelayout",
		defaultDateLayout,
		"output dates in this layout",
	)
	flags.StringVar(
		&params.from,
		"from",
		"",
		fmt.Sprintf("product/service of input file; one of %v", fromServices),
	)
	flags.BoolVar(
		&params.asc,
		"asc",
		true,
		"sort by date in ascending order?",
	)

	// The caller should set the flags.Usage field.

	return flags, params
}
