package helm

import (
	// "os/exec"
	"path"

	"github.com/trevex/kale/pkg/stage"
	"go.starlark.net/starlark"
)

func depBuild(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	chartDir := ""
	verify := false
	if err := starlark.UnpackArgs("dep_build", args, kwargs, "chart_dir", &chartDir, "verify?", &verify); err != nil {
		return nil, err
	}
	// If directory is not absolute prepend the current project directory
	if !path.IsAbs(chartDir) {
		chartDir = path.Join(stage.Current.ProjectDir, chartDir)
	}

	return starlark.None, nil
}
