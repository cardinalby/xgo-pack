package bundle

import "path/filepath"

type Path string

func (p Path) GetContentsPath() string {
	return filepath.Join(string(p), "Contents")
}

func (p Path) GetBinDir() string {
	return filepath.Join(p.GetContentsPath(), "MacOS")
}

func (p Path) GetPlistPath() string {
	return filepath.Join(p.GetContentsPath(), "Info.plist")
}

func (p Path) GetResourcesPath() string {
	return filepath.Join(p.GetContentsPath(), "Resources")
}
