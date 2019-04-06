package util

import (
	"go.starlark.net/starlark"
)

type StarlarkFunction func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error)
