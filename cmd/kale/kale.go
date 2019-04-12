package main

import (
	"fmt"
	"os"

	"github.com/trevex/kale/pkg/cmd/kale"
)

func main() {
	cmd, err := kale.Run(os.Stdout, os.Args)
	if err != nil { // kalefile failed to execute or misbehaved
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err = cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
