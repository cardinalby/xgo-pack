package go_src

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func FindGoModFile(rootPath string) (goModFilePath string, err error) {
	var goModFound = errors.New("main pkg found")
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if d.Name() == vendorDirName || d.Name() == gitDirName {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Base(path) == goModName {
			goModFilePath = path
			return goModFound
		}
		return nil
	})
	if errors.Is(err, goModFound) {
		return goModFilePath, nil
	}
	if err != nil {
		return "", err
	}
	return "", fmt.Errorf("%s not found in '%s'", goModName, rootPath)
}

func GetModuleName(goModFilePath string) (string, error) {
	data, err := os.ReadFile(goModFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading '%s': %w", goModFilePath, err)
	}
	return modfile.ModulePath(data), nil
}
