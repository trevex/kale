package util

import (
	"crypto/sha256"
	"encoding/binary"
	"os"
	"path/filepath"
	"regexp"
)

func FileExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) { // TODO: we should check other unexpected errors as well!?
		return false
	}
	return !info.IsDir()
}

var dotDirRegex = regexp.MustCompile(`[/\\]\..*`)

func PathContainsDotDir(path string) bool {
	match := dotDirRegex.FindString(path)
	return match != ""
}

func DirChecksum(dir string) ([]byte, error) {
	h := sha256.New()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if PathContainsDotDir(path) {
			return nil
		}
		h := sha256.New()
		h.Write([]byte(path))
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(info.ModTime().UnixNano()))
		h.Write(b)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
