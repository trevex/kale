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

	"github.com/trevex/kale/pkg/project"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

func depBuild(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	chartDir := ""
	verify := false
	if err := starlark.UnpackArgs("dep_build", args, kwargs, "chart_dir", &chartDir, "verify?", &verify); err != nil {
		return nil, err
	}
	// If directory is not absolute prepend the current project directory
	if !path.IsAbs(chartDir) {
		chartDir = path.Join(project.ActiveProject.Dir, chartDir)
	}
	fmt.Println(chartDir) // DEBUG
	requirementsChecksum, err := util.FileChecksum(path.Join(chartDir, "requirements.yaml"), path.Join(chartDir, "requirements.lock"))
	if err != nil {
		return nil, err
	}
	cacheDir := path.Join(project.GetCurrentCacheDir(), fmt.Sprintf("dep_build-%s", requirementsChecksum))
	fmt.Println(cacheDir) // DEBUG
	return starlark.None, nil
}
