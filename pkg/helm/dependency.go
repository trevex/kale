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
	"github.com/trevex/kale/pkg/report"
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
		// Prepare parameter string for checksums
		paramsStr := fmt.Sprintf("%s\n%v", chartDir, verify)
		// Stage 1: build the helm dependencies
		// Calculate checksum of requirements files and current params
		c1 := cache.NewChecksumBuilder()
		c1.String(paramsStr)
		if err := c1.File(path.Join(chartDir, "requirements.yaml"), path.Join(chartDir, "requirements.lock")); err != nil {
			return nil, err
		}
		s1 := cache.NewStage("helm_dep_build1", c1.Build())
		if !s1.Exists() {
			report.Infof("Copying 'chart_dir' to '%s'...", s1.SubDir())
			util.CopyDir(chartDir, s1.Dir)
			report.Infof("Executing 'helm dep build'...")
			stdout, stderr, err := util.Exec(&util.ExecOptions{Dir: s1.Dir}, "helm", "dep", "build")
			if err != nil {
				report.Errorf("%s\n%s", stdout, stderr)
				return nil, err
			}
			report.Infof("%s", stdout)
		} else {
			report.SkipStepf("Rebuilding of helm dependencies not necessary, cache exists: %s", s1.SubDir())
		}
		// Stage 2: copy charts dir if something changed (this should save time rebuilding the dependendencies if not necessary)
		c2 := cache.NewChecksumBuilder()
		c2.String(paramsStr)
		if err := c2.Dir(chartDir); err != nil {
			return nil, err
		}
		s2 := cache.NewStage("helm_dep_build2", c2.Build())
		if !s2.Exists() {
			report.Infof("Copying 'chart_dir' to '%s'...", s2.SubDir())
			util.CopyDir(chartDir, s2.Dir)
		} else {
			report.SkipStepf("Rebuilding 'chart_dir' not necessary, cache exists: %s", s2.SubDir())
		}
		report.Infof("Copying charts from '%s' to '%s'...", s1.SubDir(), s2.SubDir())
		util.CopyDir(path.Join(s1.Dir, "charts"), path.Join(s2.Dir, "charts"))
		return starlark.String(s2.Dir), nil
	}
}
