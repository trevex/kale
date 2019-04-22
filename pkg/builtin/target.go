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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trevex/kale/pkg/stage"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

type Target struct {
	Name string
	Cmd  *cobra.Command
}

type paramCheckFunc func() (*starlark.Dict, error)

func RegisterTarget(project *Project) starlark.Value {
	return starlark.NewBuiltin("target", func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var err error
		targetFunc := &starlark.Function{}
		paramsSchema := starlark.NewDict(16)
		if err = starlark.UnpackArgs("target", args, kwargs, "func", &targetFunc, "params?", &paramsSchema); err != nil {
			return nil, err
		}
		// Create info struct
		var checkParams paramCheckFunc
		targetName := targetFunc.Name()
		target := &Target{
			Name: targetName,
			Cmd: &cobra.Command{
				Use:   fmt.Sprintf("%s [flags]", targetName),
				Short: fmt.Sprintf("Executing the '%s'-target. Parameters can be provided by env-variables, a config-file or the commandline-flags below.", targetName),
				Long:  ``,
				RunE: func(_ *cobra.Command, args []string) error {
					finalParams, err := checkParams()
					if err != nil {
						return err
					}
					// TODO: calculate checksum
					checksum, err := util.DirChecksum(stage.Current.ProjectDir) // TODO: use project?
					if err != nil {
						return err
					}
					fmt.Printf("%x\n", checksum)
					// Construct kwargs
					targetKwargs := []starlark.Tuple{starlark.Tuple{starlark.String("params"), finalParams}}
					_, err = starlark.Call(thread, targetFunc, starlark.Tuple{}, targetKwargs)
					return err
				},
			},
		}
		// Crete parameter checking function
		checkParams, err = createCheckParamsFunc(target.Cmd, paramsSchema)
		if err != nil {
			return nil, err
		}
		project.AddTarget(target)
		return starlark.None, nil

	})
}

func createCheckParamsFunc(cmd *cobra.Command, paramsSchema *starlark.Dict) (paramCheckFunc, error) {
	paramFuncs := []func() (starlark.Value, starlark.Value, error){}
	flags := cmd.Flags()
	for _, tuple := range paramsSchema.Items() {
		if tuple.Len() != 2 {
			return nil, fmt.Errorf("While iterating over parameter schema a tuple without length 2 was encountered!")
		}
		key, ok := starlark.AsString(tuple[0])
		if !ok {
			return nil, fmt.Errorf("Expected string as target parameter key!")
		}
		schema, ok := tuple[1].(*starlark.Dict)
		if !ok {
			return nil, fmt.Errorf("Expected dict as target paramter value!")
		}
		obj, err := SchemaObjectFromDict(schema)
		if err != nil {
			return nil, err
		}
		var paramFunc func() (starlark.Value, starlark.Value, error)
		switch obj.Type {
		case "string":
			str, ok := obj.Default.(string)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.StringVar(&str, key, str, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.String(str), nil
			}
		case "filename":
			filename, ok := obj.Default.(string)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.StringVar(&filename, key, filename, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				// TODO: check if file exists
				return starlark.String(key), starlark.String(filename), nil
			}
		case "bool":
			b, ok := obj.Default.(bool)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.BoolVar(&b, key, b, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.Bool(b), nil
			}
		case "int":
			i, ok := obj.Default.(int)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.IntVar(&i, key, i, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.MakeInt(i), nil
			}
		case "float":
			f, ok := obj.Default.(float64)
			if !ok {
				return nil, fmt.Errorf("Provided default '%v' not of type '%s!", obj.Default, obj.Type)
			}
			flags.Float64Var(&f, key, f, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.Float(f), nil
			}
		default:
			return nil, fmt.Errorf("Type %s not implemented as target parameter", obj.Type)
		}
		paramFuncs = append(paramFuncs, paramFunc)
	}
	return func() (*starlark.Dict, error) {
		dict := &starlark.Dict{}
		for _, f := range paramFuncs {
			k, v, err := f()
			if err != nil {
				return nil, err
			}
			dict.SetKey(k, v)
		}
		return dict, nil
	}, nil
}
