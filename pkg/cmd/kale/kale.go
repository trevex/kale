/*
Copyright 2019 The Kale Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kale

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/builtin"
	"github.com/trevex/kale/pkg/engine"
	"github.com/trevex/kale/pkg/helm"
	"github.com/trevex/kale/pkg/kubectl"
	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/project"
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
	// Setup root project as unnamed project, has to be set by kalefile
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	proj := project.New("", wd, cmd)
	proj.Activate() // Setup as current project
	// Setup modules
	mgr := module.NewManager()
	mgr.Set("kubectl", kubectl.Builder)
	mgr.Set("helm", helm.Builder)
	// Create starlark engine
	eng := engine.New(stdout)
	eng.Declare(starlark.StringDict{
		"schema":        builtin.SchemaModule(),
		"global_params": builtin.GlobalParams(flags),
		"require":       builtin.RequireModule(mgr),
		"target":        builtin.RegisterTarget(proj),
		"project":       builtin.NameProject(proj),
		"var":           builtin.VarModule(proj),
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
			proj.Dir = dir
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
