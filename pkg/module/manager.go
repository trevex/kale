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
	"fmt"

	"go.starlark.net/starlark"
)

type ModuleBuilderFunction func(*starlark.Dict) (starlark.Value, error)

type moduleMap map[string]ModuleBuilderFunction

type Manager struct {
	modules moduleMap
}

func NewManager() *Manager {
	return &Manager{
		modules: moduleMap{},
	}
}

func (m *Manager) Set(name string, builder ModuleBuilderFunction) error {
	if builder == nil {
		return fmt.Errorf("Nil-pointer provided as builder function for '%s'!", name)
	}
	m.modules[name] = builder
	return nil
}

func (m *Manager) Get(name string) (ModuleBuilderFunction, bool) {
	builder, ok := m.modules[name]
	return builder, ok
}

func (m *Manager) AddBuiltin() {}
