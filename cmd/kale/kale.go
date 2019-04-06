package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/builtin"
	"github.com/trevex/kale/pkg/engine"
	"github.com/trevex/kale/pkg/global"
	"github.com/trevex/kale/pkg/module"
	"go.starlark.net/starlark"
)

func main() {
	cmd := &cobra.Command{
		Use:          "kale [flags] [target]",
		SilenceUsage: true,
		Short:        "",
		Long:         ``,
	}
	// Persistent flags
	flags := cmd.PersistentFlags()
	flags.BoolVar(&global.DryRun, "dry-run", global.DryRun, "Whether to run target without introducing changes (default false)")
	flags.StringVarP(&global.Namespace, "namespace", "n", global.Namespace, "Which namespace the target should operate in (default \"\")")
	flags.StringVarP(&global.Release, "release", "r", global.Release, "How the artifacts of the target should be named (default \"\")")
	var kalefile string
	flags.StringVarP(&kalefile, "file", "f", "./kalefile", "Kale-compliant starlark file containing target definitions")
	var config string
	flags.StringVarP(&config, "config", "c", "", "YAML-file containing target-parameters (default \"\")")
	// Parse flags once beforehand to make them available via global variables
	cmd.ParseFlags(os.Args)
	// Setup modules
	mgr := module.NewManager()
	mgr.Set("foo", func(params *starlark.Dict) (starlark.Value, error) {
		foo := &module.Module{}
		foo.SetKeyFunc(starlark.String("sayHello"), func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			fmt.Println("Hello, go")
			return starlark.String("Hello, starlark"), nil
		})
		return foo, nil
	})
	// Create starlark engine
	eng := engine.New()
	eng.Declare(starlark.StringDict{
		"greeting": starlark.String("hello"),
		"require":  builtin.Require(mgr),
		"target":   builtin.Target(cmd),
	})
	// Load file and execute it
	err := eng.ExecFile(kalefile)
	if err != nil {
		log.Fatal(err)
	}
	eng.PrintGlobalScope()
	// Start REPL if no target was supplied
	cmd.RunE = func(_ *cobra.Command, args []string) error {
		eng.REPL()
		return nil
	}
	// Run the command
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
