/*
Copyright 2019 The Kale Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package builtin

import (
	"fmt"

	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

var (
	defaultString string  = ""
	defaultBool   bool    = false
	defaultInt    int     = 0
	defaultFloat  float64 = 0.0
)

func SchemaModule() starlark.Value {
	mod := &module.Module{}
	mod.SetKeyFunc(starlark.String("filename"), primitiveSchemaObjectFunc("filename", defaultString))
	mod.SetKeyFunc(starlark.String("string"), primitiveSchemaObjectFunc("string", defaultString))
	mod.SetKeyFunc(starlark.String("bool"), primitiveSchemaObjectFunc("bool", defaultBool))
	mod.SetKeyFunc(starlark.String("int"), primitiveSchemaObjectFunc("int", defaultInt))
	mod.SetKeyFunc(starlark.String("float"), primitiveSchemaObjectFunc("float", defaultFloat))
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

func SchemaObjectFromDict(dict *starlark.Dict) (*SchemaObject, error) {
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

func primitiveSchemaObjectFunc(name string, defaultValue interface{}) func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

		obj := NewSchemaObject(name, defaultValue)
		if err := obj.UnpackFromArgs(args, kwargs); err != nil {
			return nil, err
		}
		return obj.ToDict()
	}
}

// func schemaList(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// func schemaDict(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

func getDefaultValue(name string) (interface{}, error) {
	switch name {
	case "string", "filename":
		return defaultString, nil
	case "bool":
		return defaultBool, nil
	case "float":
		return defaultFloat, nil
	case "int":
		return defaultInt, nil
	default:
		return "", fmt.Errorf("Unknown type: %s", name)
	}
}
