package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/global"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "kale [flags] [target]",
		// SilenceUsage: true,
		Short: "",
		Long:  ``,
	}
	// Persistent flags
	flags := rootCmd.PersistentFlags()
	flags.BoolVar(&global.DryRun, "dry-run", global.DryRun, "Whether to run target without introducing changes (default false)")
	flags.StringVar(&global.Namespace, "namespace", global.Namespace, "Which namespace the target should operate in (default \"\")")
	flags.StringVar(&global.Release, "release", global.Release, "How the artifacts of the target should be named (default \"\")")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
