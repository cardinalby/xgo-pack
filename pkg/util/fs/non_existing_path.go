package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	rndutil "github.com/cardinalby/xgo-pack/pkg/util/rnd"
)

func GetNonExisting(parentDirPath string, desiredName string) (resultingName string, err error) {
	getIndexedName := func(desiredName string, index int) string {
		if index == 0 {
			return desiredName
		}
		if ext := filepath.Ext(desiredName); ext != "" {
			return desiredName[:len(desiredName)-len(ext)] + "_" + rndutil.String(5) + ext
		}
		return desiredName + "_" + strconv.Itoa(index)
	}
	for i := 0; i < 100; i++ {
		tryName := getIndexedName(desiredName, i)
		_, err := os.Stat(filepath.Join(parentDirPath, tryName))
		if os.IsNotExist(err) {
			return tryName, nil
		}
		if err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf("error finding non-existing path in '%s' with desired name '%s'",
		parentDirPath, desiredName)
}
