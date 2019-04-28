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

package helm

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/project"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

// TODO: Yet, another idea to reduce boiler plate:
//       ```
//       helm = require('helm', { 'version': '>=2.11.0', 'defaults': { 'chart_dir': './helm' } })
//       Provide some defaults via initial module require, that will be used by all subsequent calls of module functions.

func GetVersion() (string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := exec.Command("helm", "version", "--client", "--short")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to verify installed helm version by running the following command: 'helm version --client --short' returned error '%s'", err)
	}
	version := strings.Split(strings.Replace(stdout.String(), "Client: ", "", 1), "+")[0]
	return version, nil
}

func template(proj *project.Project) util.StarlarkFunction {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var (
			chart     string
			namespace string
			release   string
			values    string
		)
		if err := starlark.UnpackArgs("template", args, kwargs, "chart", &chart, "namespace?", &namespace, "release?", &release, "values?", &values); err != nil {
			return nil, err
		}
		fmt.Printf("helm chart %s\n", chart)
		return starlark.None, nil
	}
}

var Builder = func(proj *project.Project, params *starlark.Dict) (starlark.Value, error) {
	mod := &module.Module{}
	versionConstraint := ">=0.0.0"
	if v, ok, err := params.Get(starlark.String("version")); ok && err == nil {
		if providedVersionConstraint, ok := starlark.AsString(v); ok {
			versionConstraint = providedVersionConstraint
		} else {
			return nil, fmt.Errorf("'version'-field is %s, but expected to be string!", v.Type())
		}
	}
	version, err := GetVersion()
	if err != nil {
		return nil, err
	}
	ok, err := util.CheckVersionConstraint(versionConstraint, version)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Helm version %s does not match constraint %s!", version, versionConstraint)
	}
	mod.SetKeyFunc(starlark.String("template"), template(proj))
	mod.SetKeyFunc(starlark.String("dep_build"), depBuild(proj))
	return mod, nil
}
