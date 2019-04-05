package engine

import (
	"fmt"
	"log"

	"go.starlark.net/repl"
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

type Engine struct {
	scope       starlark.StringDict
	predeclared starlark.StringDict
	thread      *starlark.Thread
}

func New() *Engine {
	return &Engine{
		predeclared: starlark.StringDict{},
		thread: &starlark.Thread{
			Name:  "kale",
			Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) }, // TODO: Load function
		},
	}
}

func (eng *Engine) Declare(declaration starlark.StringDict) {
	eng.predeclared = declaration // TODO: merge instead of overwrite!?
}

func (eng *Engine) ExecFile(filename string) error {
	var err error
	eng.scope, err = starlark.ExecFile(eng.thread, filename, data, eng.predeclared)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			log.Fatal(evalErr.Backtrace())
		}
		return err
	}
	return nil
}

func (eng *Engine) REPL() {
	eng.thread.Name = "REPL"
	repl.REPL(eng.thread, eng.scope)
}

func (eng *Engine) PrintGlobalScope() {
	fmt.Println("\nGlobals:")
	for _, name := range eng.scope.Keys() {
		v := eng.scope[name]
		fmt.Printf("%s (%s) = %s\n", name, v.Type(), v.String())
	}
}
