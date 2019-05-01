package report

import (
	"fmt"
	"io"
	"os"

	"github.com/logrusorgru/aurora"
)

var stdout io.Writer = os.Stdout
var Style aurora.Aurora = aurora.NewAurora(true)

func SetOutput(out io.Writer, useColors bool) {
	stdout = out
	Style = aurora.NewAurora(useColors)
}

func SkipStepf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, "%s [%s]\n", fmt.Sprintf(format, args...), Style.Bold(Style.Yellow("skip")))
}

func Infof(format string, args ...interface{}) {
	fmt.Fprintf(stdout, "%s\n", fmt.Sprintf(format, args...))
}
