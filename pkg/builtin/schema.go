package builtin

import (
	"github.com/trevex/kale/pkg/module"
	"go.starlark.net/starlark"
)

func SchemaModule() starlark.Value {
	mod := &module.Module{}
	mod.SetKeyFunc(starlark.String("filename"), primitiveSchemaObjectFunc("filename", ""))
	mod.SetKeyFunc(starlark.String("string"), primitiveSchemaObjectFunc("string", ""))
	mod.SetKeyFunc(starlark.String("bool"), primitiveSchemaObjectFunc("bool", false))
	mod.SetKeyFunc(starlark.String("int"), primitiveSchemaObjectFunc("int", 0))
	mod.SetKeyFunc(starlark.String("float"), primitiveSchemaObjectFunc("float", 0.0))
	return mod
}

type SchemaObject struct {
	Type     string
	Required bool
	Default  interface{}
}

func NewSchemaObject(name string, defaultValue interface{}) *SchemaObject {
	return &SchemaObject{
		Type:     name,
		Required: false,
		Default:  defaultValue,
	}
}

func (obj *SchemaObject) UnpackFromArgs(args starlark.Tuple, kwargs []starlark.Tuple) error {
	if err := starlark.UnpackArgs(obj.Type, args, kwargs, "required?", &obj.Required, "default?", &obj.Default); err != nil {
		return err
	}
	return nil
}

func (obj *SchemaObject) ToDict() starlark.Value {
	dict := starlark.NewDict(16)
	dict.SetKey(starlark.String("type"), starlark.String(obj.Type))
	dict.SetKey(starlark.String("required"), starlark.Bool(obj.Required))
	if obj.Default != nil {
		// TODO: implement helper function to convert any type to starlark value!
	}
	return dict
}

func primitiveSchemaObjectFunc(name string, defaultValue interface{}) func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

		obj := NewSchemaObject("filename", "")
		if err := obj.UnpackFromArgs(args, kwargs); err != nil {
			return nil, err
		}
		return obj.ToDict(), nil
	}
}

// func schemaList(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// func schemaDict(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
