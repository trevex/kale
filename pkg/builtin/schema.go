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
	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/schema"
	"go.starlark.net/starlark"
)

func SchemaModule() starlark.Value {
	mod := &module.Module{}
	mod.SetKeyFunc(starlark.String("filename"), primitiveSchemaObjectFunc("filename", schema.DefaultString))
	mod.SetKeyFunc(starlark.String("string"), primitiveSchemaObjectFunc("string", schema.DefaultString))
	mod.SetKeyFunc(starlark.String("bool"), primitiveSchemaObjectFunc("bool", schema.DefaultBool))
	mod.SetKeyFunc(starlark.String("int"), primitiveSchemaObjectFunc("int", schema.DefaultInt))
	mod.SetKeyFunc(starlark.String("float"), primitiveSchemaObjectFunc("float", schema.DefaultFloat))
	return mod
}

func primitiveSchemaObjectFunc(name string, defaultValue interface{}) func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

		obj := schema.New(name, defaultValue)
		if err := obj.UnpackFromArgs(args, kwargs); err != nil {
			return nil, err
		}
		return obj.ToDict()
	}
}

// func schemaList(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// func schemaDict(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
