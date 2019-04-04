package main

import (
	"fmt"
	"log"

	"go.starlark.net/starlark"
)

const data = `
print(greeting + ", world")

squares = [x*x for x in range(10)]

bar = require("foo")
baz = bar["sayHello"]()
print(baz)
`

type StarlarkModuleFunc func(*starlark.Dict) (starlark.Value, error)
type StarlarkModuleMap map[string]StarlarkModuleFunc

var modules StarlarkModuleMap

func init() {
	modules = make(StarlarkModuleMap)
	modules["foo"] = func(params *starlark.Dict) (starlark.Value, error) {
		foo := starlark.NewDict(1)
		foo.SetKey(starlark.String("sayHello"), starlark.NewBuiltin("sayHello", func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
			fmt.Println("Hello, go")
			return starlark.String("Hello, starlark"), nil
		}))
		return foo, nil
	}
}

func Require(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name string
	var params starlark.Dict
	if err := starlark.UnpackArgs("require", args, kwargs, "name", &name, "params?", &params); err != nil {
		return nil, err
	}
	if modFunc, ok := modules[name]; ok {
		if mod, err := modFunc(&params); err == nil {
			return mod, nil
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("TODO MODULE NOT FOUND ERROR")
	}
}

func main() {
	// Setup thread and environment
	thread := &starlark.Thread{
		Name:  "kale",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}
	predeclared := starlark.StringDict{
		"greeting": starlark.String("hello"),
		"require":  starlark.NewBuiltin("require", Require),
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
