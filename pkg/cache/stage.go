package cache

import (
	"fmt"
	"path"
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
