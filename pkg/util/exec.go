package util

import (
	"bytes"
	"io"
	"os/exec"
)

type ExecOptions struct {
	Stdin io.Reader
	Dir   string
}

func Exec(opts *ExecOptions, name string, args ...string) (string, string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := exec.Command(name, args...)
	if opts != nil {
		if opts.Stdin != nil {
			cmd.Stdin = opts.Stdin
		}
		if opts.Dir != "" {
			cmd.Dir = opts.Dir
		}
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run() // TODO: good way to wrap command execution and make them async by default?
	return stdout.String(), stderr.String(), err
}
