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

package util

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mr-tron/base58"
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

func DirChecksum(dir string) (string, error) {
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
		return "", err
	}
	return checksumToString(h.Sum(nil)), nil
}

func FileChecksum(filepaths ...string) (string, error) {
	h := sha256.New()
	for _, f := range filepaths {
		h.Write([]byte(fmt.Sprintf("%s/n", f)))
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return "", err
		}
		h.Write(data)
	}
	return checksumToString(h.Sum(nil)), nil
}

func checksumToString(c []byte) string {
	return string(base58.Encode(c))
}
