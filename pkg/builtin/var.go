package builtin

import (
	"github.com/spf13/pflag"
	"github.com/trevex/kale/pkg/module"
	"go.starlark.net/starlark"
)

type GlobalVariables struct {
	DryRun    bool
	Namespace string
	Release   string
}

var Var = &GlobalVariables{
	DryRun:    false,
	Namespace: "",
	Release:   "",
}

func VarModule() starlark.Value {
	mod := &module.Module{}
	mod.SetKey(starlark.String("dry_run"), starlark.Bool(Var.DryRun))
	mod.SetKey(starlark.String("namespace"), starlark.String(Var.Namespace))
	mod.SetKey(starlark.String("release"), starlark.String(Var.Release))
	return mod
}

func VarFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&Var.DryRun, "dry-run", Var.DryRun, "Whether to run target without introducing changes (default false)")
	flags.StringVarP(&Var.Namespace, "namespace", "n", Var.Namespace, "Which namespace the target should operate in (default \"\")")
	flags.StringVarP(&Var.Release, "release", "r", Var.Release, "How the artifacts of the target should be named (default \"\")")

}
