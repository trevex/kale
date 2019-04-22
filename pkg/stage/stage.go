package stage

import (
	"fmt"
	"path"
)

type Stage struct {
	Dir        string
	ProjectDir string
	Checksum   []byte
	Parent     *Stage
}

var Root *Stage = &Stage{
	Dir:        "./.kale",
	ProjectDir: "./",
	Checksum:   []byte{},
	Parent:     nil,
}

var Current *Stage = Root

func PushStage(prefix string, checksum []byte) *Stage {
	s := &Stage{
		Dir:        path.Join(Current.Dir, fmt.Sprintf("%s-%x", prefix, checksum)),
		ProjectDir: Current.ProjectDir,
		Checksum:   checksum,
		Parent:     Current,
	}
	Current = s
	return s
}

func PopStage() *Stage {
	s := Current
	Current = s.Parent
	return s
}

func (s *Stage) IsRoot() bool {
	return s.Parent == nil
}
