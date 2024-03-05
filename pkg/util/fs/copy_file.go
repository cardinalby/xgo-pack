package fsutil

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()
	return copyFileImpl(in, src, dst)
}

func CopyFsFile(fs fs.FS, src, dst string) error {
	in, err := fs.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()
	return copyFileImpl(in, "embed:"+src, dst)
}

func copyFileImpl(src fs.File, srcName string, dst string) (err error) {
	srcFileInfo, err := src.Stat()
	if err != nil {
		return fmt.Errorf("srcFileInfo erroron src: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("error creating dst dir: %w", err)
	}
	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcFileInfo.Mode())
	if err != nil {
		return fmt.Errorf("error creating '%s' file: %w", dst, err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, src); err != nil {
		return fmt.Errorf("error copying file from '%s' to '%s': %w", srcName, dst, err)
	}
	err = out.Sync()
	return
}
