package helm

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

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

func template(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
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

var Builder = func(params *starlark.Dict) (starlark.Value, error) {
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
	mod.SetKeyFunc(starlark.String("template"), template)
	mod.SetKeyFunc(starlark.String("dep_build"), depBuild)
	return mod, nil
}
