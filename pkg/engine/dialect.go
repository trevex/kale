package engine

import (
	"go.starlark.net/resolve"
)

func init() {
	resolve.AllowFloat = true
	resolve.AllowSet = true
	resolve.AllowLambda = true
	resolve.AllowNestedDef = true
	resolve.AllowRecursion = true
	resolve.AllowGlobalReassign = true
}
