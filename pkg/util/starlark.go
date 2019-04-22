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
package util

import (
	"fmt"

	"go.starlark.net/starlark"
)

type StarlarkFunction func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error)

func StarlarkAsBool(x starlark.Value) (bool, bool) { v, ok := x.(starlark.Bool); return bool(v), ok }

func StarlarkUnpackToInterface(x starlark.Value) (interface{}, error) {
	switch x.Type() {
	case "string":
		str, _ := starlark.AsString(x)
		return str, nil
	case "bool":
		b, _ := StarlarkAsBool(x)
		return b, nil
	case "float":
		f, _ := starlark.AsFloat(x)
		return f, nil
	case "int":
		i, _ := starlark.AsInt32(x)
		return i, nil
	default:
		return nil, fmt.Errorf("Type conversion not implemented for %s", x.Type())
	}
}

func ToStarlarkValue(v interface{}) (starlark.Value, error) {
	// string
	if x, ok := v.(string); ok {
		return starlark.String(x), nil
	}
	if x, ok := v.(starlark.String); ok {
		return x, nil
	}
	// bool
	if x, ok := v.(bool); ok {
		return starlark.Bool(x), nil
	}
	if x, ok := v.(starlark.Bool); ok {
		return x, nil
	}
	// float
	if x, ok := v.(float64); ok {
		return starlark.Float(x), nil
	}
	if x, ok := v.(starlark.Float); ok {
		return x, nil
	}
	// int
	if x, ok := v.(int); ok {
		return starlark.MakeInt(x), nil
	}
	if x, ok := v.(starlark.Int); ok {
		return x, nil
	}
	return nil, fmt.Errorf("Unsupported type: %v", v)

}
