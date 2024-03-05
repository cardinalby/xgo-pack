package fs_resource

import (
	"os"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/util/fs"
)

type PermanentDir struct {
	Path string
}

func NewPermanentDir(path string) *PermanentDir {
	return &PermanentDir{
		Path: path,
	}
}

func (t *PermanentDir) GetPath() string {
	return t.Path
}

type TempDir struct {
	Path string
}

func NewTempDir(parentDirPath string, dirName string) (*TempDir, error) {
	name, err := fsutil.GetNonExisting(parentDirPath, dirName)
	if err != nil {
		return nil, err
	}
	return &TempDir{
		Path: filepath.Join(parentDirPath, name),
	}, nil
}

func (t *TempDir) Dispose() error {
	return os.RemoveAll(t.Path)
}

func (t *TempDir) GetPath() string {
	return t.Path
}
