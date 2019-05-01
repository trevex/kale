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
	"github.com/trevex/kale/pkg/project"
	"go.starlark.net/starlark"
)

func RegisterTarget(proj *project.Project) starlark.Value {
	return starlark.NewBuiltin("target", func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var (
			err        error
			targetFunc starlark.Callable
		)
		paramsSchema := starlark.NewDict(16)
		if err = starlark.UnpackArgs("target", args, kwargs, "func", &targetFunc, "params?", &paramsSchema); err != nil {
			return nil, err
		}
		//
		_, err = proj.AddTarget(targetFunc.Name(), thread, targetFunc, paramsSchema)
		if err != nil {
			return nil, err
		}
		return starlark.None, nil
	})
}
