package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rafaelespinoza/csvtx/internal/cmd"
)

func main() {
	root := cmd.New()
	if err := root.Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
