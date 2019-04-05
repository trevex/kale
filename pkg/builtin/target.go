package builtin

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
)

func Target(rootCmd *cobra.Command) starlark.Value {
	return starlark.NewBuiltin("target", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var (
			targetImpl   starlark.Callable
			paramsSchema starlark.Dict
		)
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
				fmt.Println("TEST")
				return nil
			},
		}
		rootCmd.AddCommand(cmd)
		// Setup flags
		flags := cmd.Flags()
		test := false
		flags.BoolVar(&test, "test", false, "foo")
		// TODO: do something
		return starlark.None, nil

	})
}
