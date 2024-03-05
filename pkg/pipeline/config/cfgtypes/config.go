package cfgtypes

import (
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/build_go/config"
	"github.com/cardinalby/xgo-pack/pkg/consts"
)

const DefaultBuiltInIcon = "xgo-pack:default-icon"
const DefaultSchema = "https://raw.githubusercontent.com/cardinalby/xgo-pack/master/config_schema/config.schema.v1.json"

type PresetName string

type ArchTarget interface {
	ShouldKeepBin() bool
	ShouldBuildBin() bool
	GetOutDir() string
	SetOutDir(outDir string)
}

func asArchTargets[V ArchTarget](m map[consts.Arch]V) map[consts.Arch]ArchTarget {
	res := make(map[consts.Arch]ArchTarget, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

type OsTargets interface {
	GetArches() map[consts.Arch]ArchTarget
	GetCommonCfg() TargetsCommon
	SetCommonCfg(cfg TargetsCommon)
}

type Src struct {
	MainPkg string `json:"main_pkg,omitempty"`
	// Icon is a path to the icon file relative to Root. It will be converted to needed formats via imagemagick,
	// so for psd files you can use icon.psd[0] to convert all layers
	Icon string `json:"icon,omitempty"`
}

type Targets struct {
	// Common will be used as defaults for all other targets
	Common  TargetsCommon `json:"common,omitempty"`
	Windows TargetWindows `json:"windows,omitempty"`
	Macos   TargetMacos   `json:"macos,omitempty"`
	Linux   TargetLinux   `json:"linux,omitempty"`
}

func (t *Targets) GetOsCommonCfg(os consts.Os) *TargetsCommon {
	switch os {
	case consts.OsWindows:
		return &t.Windows.Common.TargetsCommon
	case consts.OsDarwin:
		return &t.Macos.Common.TargetsCommon
	case consts.OsLinux:
		return &t.Linux.Common.TargetsCommon
	default:
		panic("unknown os")
	}
}

func (t *Targets) GetOsArches() map[consts.Os]OsTargets {
	return map[consts.Os]OsTargets{
		consts.OsWindows: &t.Windows,
		consts.OsDarwin:  &t.Macos,
		consts.OsLinux:   &t.Linux,
	}
}

type XGoConfig struct {
	// Go release to use for cross compilation
	GoVersion string `json:"go_version,omitempty"`
	// Set a Global Proxy for Go Modules
	GoProxy string `json:"go_proxy,omitempty"`
	Verbose *bool  `json:"verbose,omitempty"`
}

type Config struct {
	// Used for JSON schema
	Schema string `json:"$schema,omitempty"`
	// Presets is a preset config names list that will be used as a base for the config.
	// Presets will be applied in the order of appearance in the list, so the last one will override the previous ones.
	Presets []PresetName `json:"presets,omitempty"`
	// Root path of the project. Absolute or relative to working directory.
	Root string `json:"root,omitempty"`
	// DistDir is path (relative to root) of a directory for final build artifacts.
	// If not set, "dist" dir will be used
	DistDir string `json:"dist_dir,omitempty"`
	// TmpDir is path (relative to root) of a temp directory for temporary build artifacts. If not set, a temporary
	// directory will be created in DistDir
	TmpDir string `json:"tmp_dir,omitempty"`
	// Src sets paths of source files and an icon
	Src Src `json:"src,omitempty"`
	// XGoConfig is a config for xgo tool
	XGo config.XGoConfig `json:"xgo,omitempty"`
	// Targets is a config for building targets
	Targets Targets `json:"targets,omitempty"`
}

func (c Config) GetTmpDirPath() string {
	return filepath.Join(c.Root, c.TmpDir)
}

func (c Config) GetDistDirPath() string {
	return filepath.Join(c.Root, c.DistDir)
}
