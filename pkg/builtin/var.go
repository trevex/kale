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
