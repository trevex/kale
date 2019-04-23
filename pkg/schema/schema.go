package schema

import (
	"fmt"

	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

var (
	DefaultString string  = ""
	DefaultBool   bool    = false
	DefaultInt    int     = 0
	DefaultFloat  float64 = 0.0
)

type SchemaObject struct {
	Type     string
	Required bool
	Default  interface{}
}

func New(name string, defaultValue interface{}) *SchemaObject {
	return &SchemaObject{
		Type:     name,
		Required: false,
		Default:  defaultValue,
	}
}

func FromDict(dict *starlark.Dict) (*SchemaObject, error) {
	obj := &SchemaObject{}
	// type
	v, ok, err := dict.Get(starlark.String("type"))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Field 'type' missing! Not a schema-object?")
	}
	obj.Type, ok = starlark.AsString(v)
	if !ok {
		return nil, fmt.Errorf("Field 'type' not a string! Not a schema-object?")
	}
	// required
	v, ok, err = dict.Get(starlark.String("required"))
	if err != nil {
		return nil, err
	}
	if !ok {
		obj.Required = false
	} else {
		obj.Required, ok = util.StarlarkAsBool(v)
		if !ok {
			return nil, fmt.Errorf("Field 'required' not a bool! Not a schema-object?")
		}
	}
	// default
	v, ok, err = dict.Get(starlark.String("default"))
	if err != nil {
		return nil, err
	}
	if !ok { // no default provided
		obj.Default, err = getDefaultValue(obj.Type)
		if err != nil {
			return nil, err
		}
	} else { // default provided
		obj.Default, err = util.StarlarkUnpackToInterface(v)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func (obj *SchemaObject) ToDict() (starlark.Value, error) {
	dict := starlark.NewDict(16)
	dict.SetKey(starlark.String("type"), starlark.String(obj.Type))
	dict.SetKey(starlark.String("required"), starlark.Bool(obj.Required))
	if obj.Default != nil {
		v, err := util.ToStarlarkValue(obj.Default)
		if err != nil {
			return nil, err
		}
		dict.SetKey(starlark.String("default"), v)
	}
	return dict, nil
}

func (obj *SchemaObject) UnpackFromArgs(args starlark.Tuple, kwargs []starlark.Tuple) error {
	if err := starlark.UnpackArgs(obj.Type, args, kwargs, "required?", &obj.Required, "default?", &obj.Default); err != nil {
		return err
	}
	return nil
}

func getDefaultValue(name string) (interface{}, error) {
	switch name {
	case "string", "filename":
		return DefaultString, nil
	case "bool":
		return DefaultBool, nil
	case "float":
		return DefaultFloat, nil
	case "int":
		return DefaultInt, nil
	default:
		return "", fmt.Errorf("Unknown type: %s", name)
	}
}
