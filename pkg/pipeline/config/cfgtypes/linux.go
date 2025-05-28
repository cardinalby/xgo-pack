package cfgtypes

import (
	"github.com/cardinalby/xgo-pack/pkg/consts"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

type TargetLinuxDebDesktopEntry struct {
	// Defines if desktop entry should be added to the package
	// If empty, `true` will be used
	AddDesktopEntry *bool `json:"add_desktop_entry,omitempty"`
	// Defines if icon should be added to the package
	// If empty, `true` will be used (only if add_desktop_entry is true)
	AddIcon *bool `json:"add_icon,omitempty"`
	// If empty, "/usr/share/icons/[common.identifierProductName].png" will be used
	DstIconPath string `json:"dst_icon_path,omitempty"`
	// Desktop entry name.
	// If empty, common.product_name will be used
	Name string `json:"name,omitempty"`
	// Desktop entry type. If empty, "Application" will be used
	Type string `json:"type,omitempty"`
	// Desktop entry Terminal key.
	// If empty, `true` will be used
	Terminal *bool `json:"terminal,omitempty"`
	// Desktop entry NoDisplay key.
	// If empty, `false` will be used
	NoDisplay *bool `json:"no_display,omitempty"`
	// Desktop entry mime type.
	MimeType string `json:"mime_type,omitempty"`
}

type TargetLinuxDeb struct {
	// Name of the resulting deb package file relative to arch `out_dir`
	// If empty, [common.product_name].deb will be used
	DebName string `json:"deb_name,omitempty"`
	// Path to custom nfpm config relative to the root.
	// If set, all other fields will be ignored and nfpm will be used with this config
	// and "XGO_PACK_" env variables that can be used in config fields as placeholders:
	// ${XGO_PACK_DASHED_PRODUCT_NAME}, ${XGO_PACK_ARCH}, ${XGO_PACK_VERSION}
	CustomNfpmConfig string `json:"custom_nfpm_config,omitempty"`
	// If empty, common.ProductName will be used.
	// For allowed format see https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-source
	// The string will be lower-cased and not supported symbols will be replaced with '-'
	Name string `json:"name,omitempty"`
	// If empty, "default" will be used
	Section string `json:"section,omitempty"`
	// Recommended to be filled. If empty common.identifier without last part will be used
	Maintainer  string `json:"maintainer,omitempty"`
	Description string `json:"description,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	Homepage    string `json:"homepage,omitempty"`
	License     string `json:"license,omitempty"`
	// Additional files to include to the package (local path -> destination abs path)
	Contents map[string]string `json:"contents,omitempty"`
	// Destination path for bin file in the package (as absolute path in the dest system)
	// If empty, "/usr/bin/{common.bin_name}" will be used
	DstBinPath string `json:"dst_bin_path,omitempty"`
	// desktop entry file options
	DesktopEntry TargetLinuxDebDesktopEntry `json:"desktop_entry,omitempty"`
}

type TargetLinuxArch struct {
	// OutDir is a path relative to Config.DistDir where the final artifacts will be placed
	// If not set, "linux_[arch]" will be used
	OutDir string `json:"out_dir,omitempty"`
	// Defines if binary should be built
	// If empty, the binary will be created only if it's needed for deb package
	BuildBin *bool `json:"build_bin,omitempty"`
	// Defines if deb package should be built
	// If empty, the deb package will not be created
	BuildDeb *bool `json:"build_deb,omitempty"`
}

func (t *TargetLinuxArch) GetOutDir() string {
	return t.OutDir
}

func (t *TargetLinuxArch) ShouldBuildBin() bool {
	return typeutil.PtrValueOrDefault(t.BuildBin) || typeutil.PtrValueOrDefault(t.BuildDeb)
}

func (t *TargetLinuxArch) ShouldKeepBin() bool {
	return typeutil.PtrValueOrDefault(t.BuildBin)
}

func (t *TargetLinuxArch) ShouldBuildDeb() bool {
	return typeutil.PtrValueOrDefault(t.BuildDeb)
}

func (t *TargetLinuxArch) SetOutDir(outDir string) {
	t.OutDir = outDir
}

type TargetLinuxCommon struct {
	TargetsCommon `json:",inline"`
	// Deb is a config for deb package
	Deb TargetLinuxDeb `json:"deb,omitempty"`
}

type TargetLinux struct {
	// Common will be used as defaults for all architectures
	Common TargetLinuxCommon `json:"common,omitempty"`
	Arm64  TargetLinuxArch   `json:"arm64,omitempty"`
	Amd64  TargetLinuxArch   `json:"amd64,omitempty"`
}

func (t *TargetLinux) GetLinuxArches() map[consts.Arch]*TargetLinuxArch {
	return map[consts.Arch]*TargetLinuxArch{
		consts.ArchArm64: &t.Arm64,
		consts.ArchAmd64: &t.Amd64,
	}
}

func (t *TargetLinux) GetArches() map[consts.Arch]ArchTarget {
	return asArchTargets(t.GetLinuxArches())
}

func (t *TargetLinux) GetCommonCfg() TargetsCommon {
	return t.Common.TargetsCommon
}

func (t *TargetLinux) SetCommonCfg(cfg TargetsCommon) {
	t.Common.TargetsCommon = cfg
}
