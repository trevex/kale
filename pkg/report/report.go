package report

import (
	"fmt"
	"io"
	"os"

	"github.com/logrusorgru/aurora"
)

var stdout io.Writer = os.Stdout
var au aurora.Aurora = aurora.NewAurora(true)

func SetOutput(out io.Writer, useColors bool) {
	stdout = out
	au = aurora.NewAurora(useColors)
}

func SkipStepf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, "%s %s\n", au.Bold(au.Yellow("Skipping:")), fmt.Sprintf(format, args...))
}
