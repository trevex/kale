package engine

import (
	"fmt"
	"io/ioutil"
	"log"

	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

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
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
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
