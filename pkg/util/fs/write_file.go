package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteFile(path string, data []byte) error {
	return writeFileWithPerm(path, data, 0755)
}

func WriteBinFile(path string, data []byte) error {
	return writeFileWithPerm(path, data, 0777)
}

func writeFileWithPerm(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error preparing file dir '%s': %w", dir, err)
	}
	return os.WriteFile(path, data, perm)
}
