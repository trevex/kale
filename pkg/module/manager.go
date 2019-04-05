package module

import (
	"fmt"

	"go.starlark.net/starlark"
)

type ModuleBuilderFn func(*starlark.Dict) (starlark.Value, error)

type moduleMap map[string]ModuleBuilderFn

type Manager struct {
	modules moduleMap
}

func NewManager() *Manager {
	return &Manager{
		modules: moduleMap{},
	}
}

func (m *Manager) Set(name string, builder ModuleBuilderFn) error {
	if builder == nil {
		return fmt.Errorf("Nil-pointer provided as builder function for '%s'!", name)
	}
	m.modules[name] = builder
	return nil
}

func (m *Manager) Get(name string) (ModuleBuilderFn, bool) {
	builder, ok := m.modules[name]
	return builder, ok
}

func (m *Manager) AddBuiltin() {}

func (m *Manager) RequireFn() starlark.Value {
	return starlark.NewBuiltin("require", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var params starlark.Dict
		if err := starlark.UnpackArgs("require", args, kwargs, "name", &name, "params?", &params); err != nil {
			return nil, err
		}
		if builder, ok := m.Get(name); ok {
			if mod, err := builder(&params); err == nil {
				return mod, nil
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Unable to find module '%s'!", name)
		}
	})
}
