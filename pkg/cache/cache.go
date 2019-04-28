package cache

import (
	"os"
)

var (
	segments []string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		segments = []string{"."}
	} else {
		segments = []string{wd}
	}
}

func SetRootDir(dir string) {
	segments[0] = dir
}
