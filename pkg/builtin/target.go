package builtin

import (
	"fmt"

	"go.starlark.net/starlark"
)

func target(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		targetImpl   starlark.Callable
		paramsSchema starlark.Dict
	)
	if err := starlark.UnpackArgs("target", args, kwargs, "func", &targetImpl, "params?", &paramsSchema); err != nil {
		return nil, err
	}
	fmt.Println(targetImpl.Name())
	// TODO: do something
	return starlark.None, nil
}

var Target = starlark.NewBuiltin("target", target)
