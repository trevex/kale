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

package kubectl

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/project"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v2"
)

type kubectlVersion struct {
	ClientVersion struct {
		GitVersion string `yaml:"gitVersion"`
	} `yaml:"clientVersion"`
}

func GetVersion() (string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := exec.Command("kubectl", "version", "--client", "-oyaml")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to verify installed kubectl version by running the following command: 'kubectl version --client -oyaml' returned '%s'", err)
	}
	v := kubectlVersion{}
	err = yaml.Unmarshal(stdout.Bytes(), &v)
	if err != nil {
		return "", fmt.Errorf("Unable to unmarshal kubectl version information: %v", err)
	}
	return v.ClientVersion.GitVersion, nil
}

func apply(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var input string
	var dryRun bool
	if err := starlark.UnpackArgs("apply", args, kwargs, "input", &input, "dry_run?", &dryRun); err != nil {
		return nil, err
	}
	fmt.Printf("kubectl input %s\n", input)
	return starlark.None, nil
}

var Builder = func(proj *project.Project, params *starlark.Dict) (starlark.Value, error) {
	mod := &module.Module{}
	versionConstraint := ">= 0.0.0"
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
		return nil, fmt.Errorf("Kubectl version %s does not match constraint %s!", version, versionConstraint)
	}
	mod.SetKeyFunc(starlark.String("apply"), apply)
	return mod, nil
}
