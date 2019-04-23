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

package builtin

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
)

var (
	projectNameRegex = regexp.MustCompile(`^[A-Za-z]{1}[A-Za-z0-9-]*[A-Za-z0-9]{1}$`)
)

type Project struct {
	Name         string
	Targets      []*Target
	Dependencies []*Project
	Cmd          *cobra.Command
}

func NameProject(p *Project) starlark.Value {
	return starlark.NewBuiltin("project", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		name := ""
		if err := starlark.UnpackArgs("project", args, kwargs, "name", &name); err != nil {
			return nil, err
		}
		p.Name = name
		return starlark.None, nil
	})
}

func NewProject(name string, cmd *cobra.Command) *Project {
	return &Project{
		Name:         name,
		Targets:      []*Target{},
		Dependencies: []*Project{},
		Cmd:          cmd,
	}
}

func (p *Project) AddTarget(target *Target) {
	p.Targets = append(p.Targets, target)
	p.Cmd.AddCommand(target.Cmd)
}

func (p *Project) ValidateName() error {
	if !projectNameRegex.MatchString(p.Name) {
		return fmt.Errorf("'%s' is not a valid project name! It hast to match '%s'.", p.Name, projectNameRegex.String())
	}
	return nil
}
