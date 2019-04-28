package cache

import (
	"fmt"
	"path"
)

type Stage struct {
	Dir string
}

func NewStage(prefix, checksum string) *Stage {
	return &Stage{
		Dir: path.Join(GetCacheDir(), fmt.Sprintf("%s-%s", prefix, checksum)),
	}
}
