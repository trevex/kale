package report

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

func SetOutput(out io.Writer) {
	stdout = out
}

func SkipStepf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, "SKIPPING: %s\n", fmt.Sprintf(format, args...))
}
