/*
Copyright 2019 The Kale Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package engine

import (
	"fmt"
	"io"
	"io/ioutil"

	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

type Engine struct {
	scope       starlark.StringDict
	predeclared starlark.StringDict
	thread      *starlark.Thread
}

func New(stdout io.Writer) *Engine {
	return &Engine{
		predeclared: starlark.StringDict{},
		thread: &starlark.Thread{
			Name:  "kale",
			Print: func(_ *starlark.Thread, msg string) { fmt.Fprintln(stdout, msg) }, // TODO: Load function
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
			return fmt.Errorf("%s", evalErr.Backtrace())
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
