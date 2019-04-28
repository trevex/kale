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

package helm

import (
	// "os/exec"
	"fmt"
	"path"

	"github.com/trevex/kale/pkg/cache"
	"github.com/trevex/kale/pkg/project"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

func depBuild(proj *project.Project) util.StarlarkFunction {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		chartDir := ""
		verify := false
		if err := starlark.UnpackArgs("dep_build", args, kwargs, "chart_dir", &chartDir, "verify?", &verify); err != nil {
			return nil, err
		}
		// If directory is not absolute prepend the current project directory
		if !path.IsAbs(chartDir) {
			chartDir = path.Join(proj.Dir, chartDir)
		}
		fmt.Println(chartDir) // DEBUG
		// Calculate checksum of requirements files and current params
		b := cache.NewChecksumBuilder()
		b.String(chartDir, fmt.Sprintf("%v", verify))
		if err := b.File(path.Join(chartDir, "requirements.yaml"), path.Join(chartDir, "requirements.lock")); err != nil {
			return nil, err
		}
		s1 := cache.NewStage("helm_dep_build1", b.Build())
		fmt.Println(s1.Dir) // DEBUG
		// TODO: check if stage exists already, if so skip
		// => maybe some standardized logging package necessary? e.g. report module?
		return starlark.None, nil
	}
}
