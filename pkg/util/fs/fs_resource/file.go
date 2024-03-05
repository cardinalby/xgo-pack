package fs_resource

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/util/fs"
)

type PermanentFile struct {
	Path string
}

func NewPermanentFile(path string) *PermanentFile {
	return &PermanentFile{
		Path: path,
	}
}

func (f *PermanentFile) GetPath() string {
	return f.Path
}

type TempFile struct {
	IsDisposed bool
	Path       string
}

func (t *TempFile) GetPath() string {
	return t.Path
}

func (t *TempFile) Dispose() error {
	if !t.IsDisposed {
		if _, err := os.Stat(t.Path); os.IsNotExist(err) {
			return nil
		}
		return os.Remove(t.Path)
	}
	return nil
}

func NewTempFile(dirPath string, fileName string) (*TempFile, error) {
	name, err := fsutil.GetNonExisting(dirPath, fileName)
	if err != nil {
		return nil, err
	}
	return &TempFile{
		Path: filepath.Join(dirPath, name),
	}, nil
}

func WriteTempFile(bytes []byte, dirPath string, fileName string) (f *TempFile, err error) {
	f, err = NewTempFile(dirPath, fileName)
	if err != nil {
		return nil, fmt.Errorf("error obtaining temp file name: %w", err)
	}
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}
	file, err := os.Create(f.Path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, file.Close())
		}
		if err != nil {
			err = errors.Join(err, f.Dispose())
		}
	}()
	if _, err = file.Write(bytes); err != nil {
		return nil, err
	}
	return f, nil
}
