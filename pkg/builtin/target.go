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
		paramsSchema := starlark.NewDict(16)
		if err := starlark.UnpackArgs("target", args, kwargs, "func", &targetFunc, "params?", &paramsSchema); err != nil {
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

func checkTargetsParams(cmd *cobra.Command, paramsSchema *starlark.Dict) (func() (*starlark.Dict, error), error) {
	// paramFuncs := []func() (starlark.Value, starlark.Value, error){}
	// for _, tuple := range paramsSchema.Items() {
	// 	if tuple.Len() != 2 {
	// 		return nil, fmt.Errorf("While iterating over parameter schema a tuple without length 2 was encountered!")
	// 	}
	// 	key := tuple[0]
	// 	schema := tuple[1]

	// }
	return nil, fmt.Errorf("TODO: impl")
}
