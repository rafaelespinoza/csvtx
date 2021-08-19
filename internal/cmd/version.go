package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/rafaelespinoza/alf"
	"github.com/rafaelespinoza/csvtx/internal/version"
)

func makeVersion(cmdName string) alf.Directive {
	return &alf.Command{
		Description: "metadata about the build",
		Setup: func(p flag.FlagSet) *flag.FlagSet {
			flags := flag.NewFlagSet(cmdName, flag.ExitOnError)
			flags.Usage = func() {
				fmt.Fprintf(flags.Output(), `Usage: %s %s

Description:

	Shows info about the build.
`, _Bin, cmdName)
			}
			return flags
		},
		Run: func(ctx context.Context) error {
			fmt.Printf("BranchName	%s\n", version.BranchName)
			fmt.Printf("BuildTime 	%s\n", version.BuildTime)
			fmt.Printf("CommitHash	%s\n", version.CommitHash)
			fmt.Printf("GoOSArch 	%s\n", version.GoOSArch)
			fmt.Printf("GoVersion 	%s\n", version.GoVersion)
			return nil
		},
	}
}
