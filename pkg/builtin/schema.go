package builtin

import (
	"github.com/trevex/kale/pkg/module"
	"go.starlark.net/starlark"
)

func SchemaModule() starlark.Value {
	mod := &module.Module{}
	mod.SetKeyFunc(starlark.String("filename"), filename)
	return mod
}

func filename(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	required := false
	if err := starlark.UnpackArgs("filename", args, kwargs, "required?", &required); err != nil {
		return nil, err
	}
	// TODO
	return starlark.String("foo"), nil
}
