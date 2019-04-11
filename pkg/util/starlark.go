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
