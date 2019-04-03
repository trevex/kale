package main

import (
	"fmt"
	"log"

	"go.starlark.net/starlark"
)

const data = `
print(greeting + ", world")

squares = [x*x for x in range(10)]
`

func main() {
	// Setup thread and environment
	thread := &starlark.Thread{
		Name:  "kompass",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}
	predeclared := starlark.StringDict{
		"greeting": starlark.String("hello"),
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
