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
	dict := createSchemaDict("filename")
	dict.SetKey(starlark.String("required"), starlark.Bool(required))
	return dict, nil
}

func createSchemaDict(schemaType string) *starlark.Dict {
	dict := starlark.NewDict(16)
	dict.SetKey(starlark.String("type"), starlark.String(schemaType))
	return dict
}
