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
