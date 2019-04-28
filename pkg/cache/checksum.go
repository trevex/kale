package cache

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"os"
	"path/filepath"

	"github.com/mr-tron/base58"
	"github.com/trevex/kale/pkg/util"
)

type ChecksumBuilder struct {
	h hash.Hash
}

func NewChecksumBuilder() *ChecksumBuilder {
	return &ChecksumBuilder{
		h: sha256.New(),
	}
}

func (b *ChecksumBuilder) Dir(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil // TODO: symlinks?
		}
		if util.PathContainsDotDir(path) {
			return nil
		}
		b.h.Write([]byte(path))
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(info.ModTime().UnixNano()))
		b.h.Write(buf)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *ChecksumBuilder) File(filepaths ...string) error {
	for _, f := range filepaths {
		b.h.Write([]byte(f))
		info, err := os.Stat(f)
		if err != nil {
			return err
		}
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(info.ModTime().UnixNano()))
		b.h.Write(buf)
	}
	return nil
}

func (b *ChecksumBuilder) Build() string {
	return string(base58.Encode(b.h.Sum(nil)))
}
