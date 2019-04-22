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
package module

import (
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

type Module struct {
	starlark.Dict
}

func (m *Module) Type() string { return "module" }

func (m *Module) Attr(name string) (starlark.Value, error) {
	if v, found, _ := m.Get(starlark.String(name)); found {
		return v, nil
	} else {
		v, err := m.Dict.Attr(name)
		return v, err
	}
}
func (m *Module) AttrNames() []string {
	names := m.Dict.AttrNames()
	for _, v := range m.Keys() {
		if str, ok := starlark.AsString(v); ok {
			names = append(names, str)
		}
	}
	return names
}

func (m *Module) SetKeyFunc(name starlark.String, fn util.StarlarkFunction) error {
	return m.SetKey(name, starlark.NewBuiltin(string(name), fn))
}
