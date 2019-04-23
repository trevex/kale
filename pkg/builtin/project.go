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

func NameProject(proj *project.Project) starlark.Value {
	return starlark.NewBuiltin("project", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		name := ""
		if err := starlark.UnpackArgs("project", args, kwargs, "name", &name); err != nil {
			return nil, err
		}
		proj.Name = name
		return starlark.None, nil
	})
}
