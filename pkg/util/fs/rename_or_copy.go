package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/util/logging"
)

func RenameOrCopyFile(src, dst string, logger logging.Logger) error {
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("error creating dst dir: %w", err)
	}
	err := os.Rename(src, dst)
	if errors.Is(err, os.ErrPermission) {
		logger.Printf("rename failed with permission error, trying to copy")
		if err := CopyFile(src, dst); err != nil {
			return fmt.Errorf("error copying file: %w", err)
		}
		return nil
	}
	return err
}
