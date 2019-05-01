package cache

import (
	"os"
	"path"
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

func GetRootDir() string {
	return segments[0]
}

func Push(subdir string) int {
	segments = append(segments, subdir)
	return len(segments)
}

func Pop(length int) {
	if length < 2 {
		return // leave root in-tact
	}
	segments = segments[:length]
}

func GetCacheDir() string {
	return path.Join(segments...)
}
