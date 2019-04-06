package kubectl

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v2"
)

type kubectlVersion struct {
	ClientVersion struct {
		GitVersion string `yaml:"gitVersion"`
	} `yaml:"clientVersion"`
}

type kubectlParams struct {
	VersionConstraint string
}

func newKubectlParams() *kubectlParams {
	return &kubectlParams{
		VersionConstraint: "v0.0.0",
	}
}

func GetVersion() (string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := exec.Command("kubectl", "version", "--client", "-oyaml")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to verify installed kubectl version by running the following command: 'kubectl version --client -oyaml'")
	}
	v := kubectlVersion{}
	err = yaml.Unmarshal(stdout.Bytes(), &v)
	if err != nil {
		return "", fmt.Errorf("Unable to unmarshal kubectl version information: %v", err)
	}
	return v.ClientVersion.GitVersion, nil
}

func apply(params *kubectlParams) util.StarlarkFunction {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var input string
		var dryRun bool
		if err := starlark.UnpackArgs("apply", args, kwargs, "input", &input, "dry_run?", &dryRun); err != nil {
			return nil, err
		}
		fmt.Printf("kubectl input %s\n", input)
		return starlark.String(params.VersionConstraint), nil
	}
}

var Builder = func(params *starlark.Dict) (starlark.Value, error) {
	mod := &module.Module{}
	parsed := newKubectlParams()
	if v, ok, err := params.Get(starlark.String("version")); ok && err == nil {
		if versionConstraint, ok := starlark.AsString(v); ok {
			parsed.VersionConstraint = versionConstraint
		} else {
			return nil, fmt.Errorf("'version'-field is %s, but expected to be string!", v.Type())
		}
	}
	version, err := GetVersion()
	if err != nil {
		return nil, err
	}
	ok, err := util.CheckVersionConstraint(parsed.VersionConstraint, version)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Kubectl version %s does not match constraint %s!", version, parsed.VersionConstraint)
	}
	mod.SetKeyFunc(starlark.String("apply"), apply(parsed))
	return mod, nil
}
