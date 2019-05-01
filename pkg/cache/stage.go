package cache

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type Stage struct {
	Name string
	Dir  string
}

func NewStage(prefix, checksum string) *Stage {
	name := fmt.Sprintf("%s-%s", prefix, checksum)
	return &Stage{
		Name: name,
		Dir:  path.Join(GetCacheDir(), name),
	}
}

func (s *Stage) SubDir() string {
	parts := strings.SplitN(s.Dir, GetRootDir(), 2)
	if len(parts) < 2 {
		return s.Name
	}
	return strings.TrimLeft(parts[1], "/\\")
}

func (s *Stage) Exists() bool {
	if _, err := os.Stat(s.Dir); !os.IsNotExist(err) {
		return true
	}
	return false
}
