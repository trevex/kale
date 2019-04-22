package helm

import (
	"go.starlark.net/starlark"
)

func depBuild(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	chartDir := ""
	verify := false
	if err := starlark.UnpackArgs("dep_build", args, kwargs, "chart_dir", &chartDir, "verify?", &verify); err != nil {
		return nil, err
	}
	return starlark.None, nil
}
