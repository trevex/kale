package builtin

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
)

type Target struct {
	Name string
	Cmd  *cobra.Command
	// TODO: schema
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
				Short: "TODO",
				Long:  `TODO`,
				RunE: func(_ *cobra.Command, args []string) error {
					finalParams, err := checkParams()
					if err != nil {
						return err
					}
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
			str := obj.Default.(string)
			flags.StringVar(&str, key, str, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.String(str), nil
			}
		case "filename":
			filename := obj.Default.(string)
			flags.StringVar(&filename, key, filename, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				// TODO: check if file exists
				return starlark.String(key), starlark.String(filename), nil
			}
		case "bool":
			b := obj.Default.(bool)
			flags.BoolVar(&b, key, b, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.Bool(b), nil
			}
		case "int":
			i := obj.Default.(int)
			flags.IntVar(&i, key, i, "TODO")
			paramFunc = func() (starlark.Value, starlark.Value, error) {
				return starlark.String(key), starlark.MakeInt(i), nil
			}
		case "float":
			f := obj.Default.(float64)
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
