package kale

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/builtin"
	"github.com/trevex/kale/pkg/engine"
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
	cmd.SetArgs(args)
	// Persistent flags
	flags := cmd.PersistentFlags()
	var (
		kalefile string
		config   string
	)
	flags.StringVarP(&kalefile, "file", "f", "./kalefile", "Kale-compliant starlark file containing target definitions")
	flags.StringVarP(&config, "config", "c", "", "YAML-file containing target-parameters (default \"\")")
	// Add builtin global flags
	builtin.VarFlags(flags)
	// Parse flags once beforehand to make them available via global variables
	cmd.ParseFlags(args)
	// TODO: if err := builtin.VarCheckRequired(cfg); err != nil {
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
		"var":     builtin.VarModule(),
	})
	// Load file and execute it
	if err := eng.ExecFile(kalefile); err != nil {
		return nil, err
	}
	// Check whether project name was set
	if err := proj.ValidateName(); err != nil {
		return nil, err
	}
	// Start REPL if no target was supplied
	cmd.RunE = func(_ *cobra.Command, args []string) error {
		eng.REPL()
		return nil
	}
	// Return the command
	return cmd, nil
}
