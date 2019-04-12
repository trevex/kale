package kale

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/builtin"
	"github.com/trevex/kale/pkg/engine"
	"github.com/trevex/kale/pkg/global"
	"github.com/trevex/kale/pkg/kubectl"
	"github.com/trevex/kale/pkg/module"
	"go.starlark.net/starlark"
)

func Run(stdout io.Writer, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "kale [flags] [target]",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "",
		Long:          ``,
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
	cmd.ParseFlags(args)
	// Setup modules
	mgr := module.NewManager()
	mgr.Set("kubectl", kubectl.Builder)
	// Setup root project
	proj := builtin.NewProject("", cmd)
	// Create starlark engine
	eng := engine.New(stdout)
	eng.Declare(starlark.StringDict{
		"require": builtin.RequireModule(mgr),
		"target":  builtin.RegisterTarget(proj),
		"project": builtin.NameProject(proj),
		"schema":  builtin.SchemaModule(),
	})
	// Load file and execute it
	err := eng.ExecFile(kalefile)
	if err != nil {
		return nil, err
	}
	// Check whether project name was set
	if proj.Name == "" {
		return nil, fmt.Errorf("No project name was provided by file: %s", kalefile)
	}
	// Start REPL if no target was supplied
	cmd.RunE = func(_ *cobra.Command, args []string) error {
		eng.REPL()
		return nil
	}
	// Return the command
	return cmd, nil
}
