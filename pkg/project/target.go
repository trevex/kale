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

package project

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/module"
	"go.starlark.net/starlark"
)

var (
	ActiveTarget *Target = nil
)

type Target struct {
	Name        string
	Cmd         *cobra.Command // TODO: maybe private?
	CheckParams func(*starlark.Dict) error
	CacheDir    string
}

func newTarget(proj *Project, name string, thread *starlark.Thread, targetFunc *starlark.Function) *Target {
	target := &Target{
		Name: name,
	}
	target.Cmd = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", target.Name),
		Short: fmt.Sprintf("Executing the '%s'-target. Parameters can be provided by env-variables, a config-file or the commandline-flags below.", target.Name),
		Long:  ``,
		RunE: func(_ *cobra.Command, args []string) error {
			// Make sure project is set as current and setup target as well,
			// but keep reference to previous target
			proj.Activate()
			prev := target.Activate()
			defer func() {
				// Deactivate and if there was a previous target, activate it
				target.Deactivate()
				if prev != nil {
					prev.Activate()
				}
			}()
			//
			fmt.Println(target.CacheDir)
			//
			params := &module.Module{} // Allows access via dot notation
			err := target.CheckParams(&params.Dict)
			if err != nil {
				return err
			}
			// Construct kwargs
			targetKwargs := []starlark.Tuple{starlark.Tuple{starlark.String("params"), params}}
			_, err = starlark.Call(thread, targetFunc, starlark.Tuple{}, targetKwargs)
			return err
		},
	}
	return target
}

func (t *Target) Activate() *Target {
	prev := ActiveTarget
	ActiveTarget = t
	return prev
}

func (t *Target) Deactivate() {
	ActiveTarget = nil
}
