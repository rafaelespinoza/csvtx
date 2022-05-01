package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rafaelespinoza/alf"
	"github.com/rafaelespinoza/csvtx/internal/convert"
	"github.com/rafaelespinoza/csvtx/internal/entity"
)

func makeNormalizeJSON(parentName, name string) alf.Directive {
	var params *normalizeParams

	out := alf.Command{
		Description: "read CSV from stdin and print normalized data to stdout as JSON",
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

				for _, dat := range data {
					if err = writeJSON(dat, out); err != nil {
						return
					}
				}
			case serviceMint:
				var data []*entity.Mint
				data, err = convert.ReadParseMint(in, params.asc)
				if err != nil {
					return
				}

				for _, dat := range data {
					if err = writeJSON(dat, out); err != nil {
						return
					}
				}
			case serviceVenmo:
				var data []*entity.Venmo
				data, err = convert.ReadParseVenmo(in, params.asc)
				if err != nil {
					return
				}

				for _, dat := range data {
					if err = writeJSON(dat, out); err != nil {
						return
					}
				}
			case serviceWellsFargo:
				var data []*entity.WellsFargo
				data, err = convert.ReadParseWellsFargo(in, params.asc)
				if err != nil {
					return
				}

				for _, dat := range data {
					if err = writeJSON(dat, out); err != nil {
						return
					}
				}
			default:
				err = fmt.Errorf("unknown source service %q", params.from)
			}

			return
		},
	}

	return &out
}

func writeJSON(data interface{}, out io.Writer) (err error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(out, "%s\n", raw)
	return
}
