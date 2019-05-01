package report

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

var stdout io.Writer = os.Stdout
var Style aurora.Aurora = aurora.NewAurora(true)

func SetOutput(out io.Writer, useColors bool) {
	stdout = out
	Style = aurora.NewAurora(useColors)
}

func sprintf(format string, args ...interface{}) string {
	return strings.TrimRight(fmt.Sprintf(format, args...), "\n")
}

func SkipStepf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, "%s (%s)\n", sprintf(format, args...), Style.Bold(Style.Yellow("skipped")))
}

func Infof(format string, args ...interface{}) {
	fmt.Fprintf(stdout, "%s\n", sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) { // TODO: stderr?
	fmt.Fprintf(stdout, "%s %s\n", Style.Bold(Style.Red("Error:")), sprintf(format, args...))
}
