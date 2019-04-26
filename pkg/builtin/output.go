package builtin

import (
	"fmt"

	"github.com/trevex/kale/pkg/project"
	"go.starlark.net/starlark"
)

func Output(proj *project.Project) starlark.Value { // TODO: needs project?
	return starlark.NewBuiltin("output", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var output starlark.Value
		if err := starlark.UnpackArgs("output", args, kwargs, "output", &output); err != nil {
			return nil, err
		}
		fmt.Println(output.String())
		return starlark.None, nil
	})
}
