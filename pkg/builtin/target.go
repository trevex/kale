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

func RegisterTarget(project *Project) starlark.Value {
	return starlark.NewBuiltin("target", func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		targetFunc := &starlark.Function{}
		paramsDict := starlark.NewDict(16)
		if err := starlark.UnpackArgs("target", args, kwargs, "func", &targetFunc, "params?", &paramsDict); err != nil {
			return nil, err
		}
		// Create info struct
		targetName := targetFunc.Name()
		target := &Target{
			Name: targetName,
			Cmd: &cobra.Command{
				Use:   fmt.Sprintf("%s [flags]", targetName),
				Short: "TODO",
				Long:  `TODO`,
				RunE: func(_ *cobra.Command, args []string) error {
					targetKwargs := []starlark.Tuple{starlark.Tuple{starlark.String("params"), starlark.String("hello, params")}}
					// TODO: proper params
					_, err := starlark.Call(thread, targetFunc, starlark.Tuple{}, targetKwargs)
					return err
				},
			},
		}
		// Setup flags
		flags := target.Cmd.Flags()
		test := false
		// TODO: proper flags
		flags.BoolVar(&test, "test", false, "foo")
		// TODO: check config?
		project.AddTarget(target)
		return starlark.None, nil

	})
}
