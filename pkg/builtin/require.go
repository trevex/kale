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
	"go.starlark.net/starlark"
)

func RequireModule(mgr *module.Manager) starlark.Value {
	return starlark.NewBuiltin("require", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		name := ""
		params := starlark.NewDict(16)
		if err := starlark.UnpackArgs("require", args, kwargs, "name", &name, "params?", &params); err != nil {
			return nil, err
		}
		if builder, ok := mgr.Get(name); ok {
			if mod, err := builder(params); err == nil {
				return mod, nil
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Unable to find module '%s'!", name)
		}
	})
}
