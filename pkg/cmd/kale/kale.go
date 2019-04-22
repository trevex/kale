package kale

import (
	"io"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/builtin"
	"github.com/trevex/kale/pkg/engine"
	"github.com/trevex/kale/pkg/helm"
	"github.com/trevex/kale/pkg/kubectl"
	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/stage"
	"github.com/trevex/kale/pkg/util"
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
	kalefile := ""
	config := ""
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
	mgr.Set("helm", helm.Builder)
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
	// Start REPL if no target was supplied
	cmd.RunE = func(_ *cobra.Command, args []string) error {
		eng.REPL()
		return nil
	}
	// Check whether kalefile exists
	if util.FileExists(kalefile) {
		// Get absolute directory of kalefile
		if dir, err := filepath.Abs(filepath.Dir(kalefile)); err != nil {
			return nil, err
		} else {
			// Update root stage directory to match kalefile
			stage.Root.ProjectDir = dir
			stage.Root.Dir = path.Join(dir, ".kale")
		}
		// Load file and execute it
		if err := eng.ExecFile(kalefile); err != nil {
			return nil, err
		}
		// Check whether project name was set
		if err := proj.ValidateName(); err != nil {
			return nil, err
		}
	}
	// Return the command
	return cmd, nil
}
