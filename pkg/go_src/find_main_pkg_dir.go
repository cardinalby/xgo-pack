package go_src

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
)

// FindMainPkgDir finds main package directory in the rootPath and returns its path relative to the rootPath.
func FindMainPkgDir(rootPath string) (string, error) {
	skipDir := ""
	var mainPkgFound = errors.New("main pkg found")
	var foundPkgDir string
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if d.Name() == vendorDirName || d.Name() == gitDirName {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Dir(path) == skipDir {
			return nil
		}
		if filepath.Ext(path) == goExt {
			pkgName, err := getGoFilePackageName(path)
			if err != nil {
				return fmt.Errorf("error getting package name from '%s': %w", path, err)
			}
			if pkgName == mainPkgName {
				foundPkgDir = filepath.Dir(path)
				return mainPkgFound
			}
			skipDir = filepath.Dir(path)
		}
		return nil
	})
	if errors.Is(err, mainPkgFound) {
		relPath, err := filepath.Rel(rootPath, foundPkgDir)
		if err != nil {
			return "", fmt.Errorf("error getting relative path of '%s' to '%s': %w", foundPkgDir, rootPath, err)
		}
		return relPath, nil
	}
	if err != nil {
		return "", err
	}
	return "", fmt.Errorf("main package not found in '%s'", rootPath)
}

func getGoFilePackageName(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, filePath, data, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("error parsing go file: %w", err)
	}
	return f.Name.Name, nil
}
