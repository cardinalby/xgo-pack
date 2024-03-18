package fsutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

var ErrIsNotRelative = fmt.Errorf("path is not relative")

func GetRelPathInsideDir(path, base string) (string, error) {
	relRes, err := filepath.Rel(base, path)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(relRes, "..") {
		return "", fmt.Errorf("%q %w to %q", path, ErrIsNotRelative, base)
	}
	return relRes, nil
}
