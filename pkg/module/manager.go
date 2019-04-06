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
