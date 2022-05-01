package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rafaelespinoza/alf"
)

// Root abstracts a top-level command from package main.
type Root interface {
	// Run is the entry point. It should be called with os.Args[1:].
	Run(ctx context.Context, args []string) error
}

// _Bin is the name of the binary file.
var _Bin = os.Args[0]

// New establishes the root command and subcommands.
func New() Root {
	deleg := alf.Delegator{
		Description: "main command for " + _Bin,
		Flags:       flag.NewFlagSet("main", flag.ExitOnError),
		Subs: map[string]alf.Directive{
			"convert":   makeConvert("convert"),
			"normalize": makeNormalize("normalize"),
			"version":   makeVersion("version"),
		},
	}

	deleg.Flags.Usage = func() {
		fmt.Fprintf(deleg.Flags.Output(), `Usage:
	%s flags

Description:

	%s is a command line tool to format financial transaction data for transfer
	between services.

Subcommands:

	These may have their own set of flags. Put them after the subcommand.

	%v

Examples:

	%s [root-flags] [subcommand] [sub-flags]

Flags:

`, _Bin, _Bin, strings.Join(deleg.DescribeSubcommands(), "\n\t"), _Bin)
		deleg.Flags.PrintDefaults()
	}

	return &alf.Root{Delegator: &deleg}
}

const (
	serviceMechanicsBank = "mechanicsbank"
	serviceMint          = "mint"
	serviceVenmo         = "venmo"
	serviceWellsFargo    = "wellsfargo"
)

var fromServices = []string{
	serviceMechanicsBank,
	serviceMint,
	serviceVenmo,
	serviceWellsFargo,
}
