package builtin

import (
	"fmt"

	"github.com/trevex/kale/pkg/project"
	"go.starlark.net/starlark"
)

func Input(proj *project.Project) starlark.Value { // TODO: needs project?
	return starlark.NewBuiltin("input", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var input starlark.Value
		if err := starlark.UnpackArgs("input", args, kwargs, "input", &input); err != nil {
			return nil, err
		}
		fmt.Printf("TODO: %s", input.String())
		return starlark.None, nil
	})
}
