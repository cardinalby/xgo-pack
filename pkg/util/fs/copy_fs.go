package fsutil

import (
	"io/fs"
	"path/filepath"
)

func CopyFs(fs fs.FS, destDir string) error {
	fileNames, err := GetFsFilenames(fs)
	if err != nil {
		return err
	}
	for _, fileName := range fileNames {
		destFile := filepath.Join(destDir, fileName)
		if err := CopyFsFile(fs, fileName, destFile); err != nil {
			return err
		}
	}
	return nil
}
