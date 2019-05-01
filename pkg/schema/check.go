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

package schema

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.starlark.net/starlark"
)

type ParameterCheck func(*starlark.Dict) error

func ConstructParameterCheck(flags *pflag.FlagSet, paramsSchema *starlark.Dict) (ParameterCheck, error) {
	paramFuncs := []func() (starlark.Value, starlark.Value, error){}
	for _, tuple := range paramsSchema.Items() {
		if tuple.Len() != 2 {
			return nil, fmt.Errorf("While iterating over parameter schema a tuple without length 2 was encountered!")
		}
		key, ok := starlark.AsString(tuple[0]) // TODO: check that key is "conform"
		if !ok {
			return nil, fmt.Errorf("Expected string as target parameter key!")
		}
		keyFlag := strings.ReplaceAll(key, "_", "-") // Use dashes in flags
		dict, ok := tuple[1].(*starlark.Dict)
		if !ok {
			return nil, fmt.Errorf("Expected dict as target paramter value!")
		}
		obj, err := FromDict(dict)
		if err != nil {
			return nil, err
		}
		var paramFunc func() (starlark.Value, starlark.Value, error)
		switch obj.Type {
		case "string":
			str, ok := obj.Default.(string)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.StringVar(&str, keyFlag, str, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.String(str), nil
			}
		case "filename":
			filename, ok := obj.Default.(string)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.StringVar(&filename, keyFlag, filename, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				// TODO: check if file exists
				return starlark.String(key), starlark.String(filename), nil
			}
		case "bool":
			b, ok := obj.Default.(bool)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.BoolVar(&b, keyFlag, b, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.Bool(b), nil
			}
		case "int":
			i, ok := obj.Default.(int)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.IntVar(&i, keyFlag, i, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.MakeInt(i), nil
			}
		case "float":
			f, ok := obj.Default.(float64)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.Float64Var(&f, keyFlag, f, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.Float(f), nil
			}
		case "input":
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				// TODO: run target func!
				return starlark.String(key), starlark.String("Hello"), nil
			}
		default:
			return nil, fmt.Errorf("Type %s not implemented as target parameter", obj.Type)
		}
		paramFuncs = append(paramFuncs, paramFunc)
	}
	return func(dict *starlark.Dict) error {
		for _, f := range paramFuncs {
			k, v, err := f()
			if err != nil {
				return err
			}
			dict.SetKey(k, v)
		}
		return nil
	}, nil
}
