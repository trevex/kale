package builtin

import (
	"fmt"

	"github.com/trevex/kale/pkg/module"
	"go.starlark.net/starlark"
)

func RequireModule(mgr *module.Manager) starlark.Value {
	return starlark.NewBuiltin("require", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		name := ""
		params := starlark.NewDict(16)
		if err := starlark.UnpackArgs("require", args, kwargs, "name", &name, "params?", &params); err != nil {
			return nil, err
		}
		if builder, ok := mgr.Get(name); ok {
			if mod, err := builder(params); err == nil {
				return mod, nil
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Unable to find module '%s'!", name)
		}
	})
}
