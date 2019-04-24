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
	"regexp"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/schema"
	"go.starlark.net/starlark"
)

var (
	projectNameRegex          = regexp.MustCompile(`^[A-Za-z]{1}[A-Za-z0-9-]*[A-Za-z0-9]{1}$`)
	ActiveProject    *Project = nil
)

type Project struct {
	Name         string
	Dir          string
	Targets      []*Target
	Dependencies []*Project
	Cmd          *cobra.Command
}

func New(name string, dir string, cmd *cobra.Command) *Project {
	return &Project{
		Name:         name,
		Dir:          dir,
		Targets:      []*Target{},
		Dependencies: []*Project{},
		Cmd:          cmd,
	}
}

func (p *Project) AddTarget(name string, thread *starlark.Thread, targetFunc *starlark.Function, paramsSchema *starlark.Dict) (*Target, error) {
	var err error
	target := newTarget(p, name, thread, targetFunc)
	target.CheckParams, err = schema.ConstructParameterCheck(target.Cmd.Flags(), paramsSchema)
	if err != nil {
		return nil, err
	}
	// Create parameter checking function
	p.Targets = append(p.Targets, target)
	p.Cmd.AddCommand(target.Cmd)
	return target, nil
}

func (p *Project) ValidateName() error {
	if !projectNameRegex.MatchString(p.Name) {
		return fmt.Errorf("'%s' is not a valid project name! It hast to match '%s'.", p.Name, projectNameRegex.String())
	}
	return nil
}

func (p *Project) Activate() {
	ActiveProject = p
}
