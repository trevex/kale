package main

import (
	"fmt"
	"log"

	"github.com/trevex/kale/pkg/builtin"
	"github.com/trevex/kale/pkg/module"
	"go.starlark.net/starlark"
)

const data = `
print(greeting + ", world")

squares = [x*x for x in range(10)]

bar = require("foo")
baz = bar.sayHello()
print(baz)

def deploy(params):
	print(params)

target(deploy)
`

func main() {
	// Setup modules
	mgr := module.NewManager()
	mgr.Set("foo", func(params *starlark.Dict) (starlark.Value, error) {
		foo := &module.Module{}
		foo.SetKeyFunc(starlark.String("sayHello"), func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			fmt.Println("Hello, go")
			return starlark.String("Hello, starlark"), nil
		})
		return foo, nil
	})
	// Setup thread and environment
	thread := &starlark.Thread{
		Name:  "kale",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}
	predeclared := starlark.StringDict{
		"greeting": starlark.String("hello"),
		"require":  builtin.Require(mgr),
		"target":   builtin.Target,
	}

	// Run the script
	globals, err := starlark.ExecFile(thread, "apparent/filename.star", data, predeclared)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			log.Fatal(evalErr.Backtrace())
		}
		log.Fatal(err)
	}

	// Print the global environment.
	fmt.Println("\nGlobals:")
	for _, name := range globals.Keys() {
		v := globals[name]
		fmt.Printf("%s (%s) = %s\n", name, v.Type(), v.String())
	}
}
