package module

import (
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

func (m *Module) SetKeyFunc(name starlark.String, fn func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error)) error {
	return m.SetKey(name, starlark.NewBuiltin(string(name), fn))
}
