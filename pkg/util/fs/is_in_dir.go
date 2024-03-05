package fsutil

import (
	"path/filepath"
	"strings"
)

func IsInDir(dir, path string) (bool, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false, err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	rel, err := filepath.Rel(absDir, absPath)
	if err != nil {
		return false, err
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)), nil
}
