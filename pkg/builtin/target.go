package builtin

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
)

func Target(rootCmd *cobra.Command) starlark.Value {
	return starlark.NewBuiltin("target", func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		targetImpl := &starlark.Function{}
		paramsSchema := starlark.NewDict(16)
		if err := starlark.UnpackArgs("target", args, kwargs, "func", &targetImpl, "params?", &paramsSchema); err != nil {
			return nil, err
		}
		targetName := targetImpl.Name()
		// Create command
		cmd := &cobra.Command{
			Use:   fmt.Sprintf("%s [flags]", targetName),
			Short: "TODO",
			Long:  `TODO`,
			RunE: func(_ *cobra.Command, args []string) error {
				targetKwargs := []starlark.Tuple{starlark.Tuple{starlark.String("params"), starlark.String("hello, params")}}
				// TODO: proper params
				_, err := starlark.Call(thread, targetImpl, starlark.Tuple{}, targetKwargs)
				return err
			},
		}
		rootCmd.AddCommand(cmd)
		// Setup flags
		flags := cmd.Flags()
		test := false
		// TODO: proper flags
		flags.BoolVar(&test, "test", false, "foo")
		// TODO: check config?
		return starlark.None, nil

	})
}
