package util

import (
	"os"
)

func FileExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) { // TODO: we should check other unexpected errors as well!?
		return false
	}
	return !info.IsDir()

}
