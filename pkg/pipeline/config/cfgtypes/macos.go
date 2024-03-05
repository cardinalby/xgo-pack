package cfgtypes

import (
	"github.com/cardinalby/xgo-pack/pkg/consts"
	"github.com/cardinalby/xgo-pack/pkg/util/maputil"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

// TargetMacosCommonBundle is an arches common settings for MacOS bundle
type TargetMacosCommonBundle struct {
	// BundleName is a name of the resulting app bundle relative to TargetMacosArch.OutDir.
	// If empty, "[targets.common.product_name].app" will be used
	BundleName string `json:"bundle_name,omitempty"`
	// HideInDock sets the plist flag in the app bundle to hide the app in dock
	HideInDock *bool `json:"hide_in_dock,omitempty"`
}

// TargetMacosCommonDmg is an arches common settings for MacOS dmg
type TargetMacosCommonDmg struct {
	// DmgName is a name of the resulting dmg file relative to TargetMacosArch.OutDir.
	// If empty, "[targets.common.product_name].dmg" will be used
	DmgName string `json:"dmg_name,omitempty"`
	// AddApplicationsSymlink is a flag to add a symlink to the /Applications folder in the dmg
	AddApplicationsSymlink *bool `json:"add_applications_symlink,omitempty"`
}

type TargetMacosCommon struct {
	TargetsCommon `json:",inline"`
	Bundle        TargetMacosCommonBundle `json:"bundle,omitempty"`
	Dmg           TargetMacosCommonDmg    `json:"dmg,omitempty"`
}

type TargetMacosArch struct {
	// OutDir is a path relative to Config.DistDir where the final artifacts will be placed
	// If not set, "macos_[arch]" will be used
	OutDir string `json:"out_dir,omitempty"`
	// Defines if binary should be built.
	// If false, the binary will be built only in case it's required to build bundle or dmg.
	// Temporary path will be used in this case.
	BuildBin *bool `json:"build_bin,omitempty"`
	// Defines if app bundle should be built.
	// If false, the bundle will be created only in case it's required to build dmg.
	// Temporary path will be used in this case.
	BuildBundle *bool `json:"build_bundle,omitempty"`
	// Defines if dmg should be built.
	// If empty, the dmg will not be created
	BuildDmg *bool `json:"build_dmg,omitempty"`
}

func (t *TargetMacosArch) ShouldBuildBin() bool {
	return typeutil.PtrValueOrDefault(t.BuildBin) || t.ShouldBuildBundle()
}

func (t *TargetMacosArch) ShouldKeepBin() bool {
	return typeutil.PtrValueOrDefault(t.BuildBin)
}

func (t *TargetMacosArch) ShouldBuildBundle() bool {
	return typeutil.PtrValueOrDefault(t.BuildBundle) || t.ShouldBuildDmg()
}

func (t *TargetMacosArch) ShouldKeepBundle() bool {
	return typeutil.PtrValueOrDefault(t.BuildBundle)
}

func (t *TargetMacosArch) ShouldBuildDmg() bool {
	return typeutil.PtrValueOrDefault(t.BuildDmg)
}

func (t *TargetMacosArch) GetOutDir() string {
	return t.OutDir
}

func (t *TargetMacosArch) SetOutDir(outDir string) {
	t.OutDir = outDir
}

type TargetMacos struct {
	// Will be used as defaults for all architectures
	Common TargetMacosCommon `json:"common,omitempty"`
	Arm64  TargetMacosArch   `json:"arm64,omitempty"`
	Amd64  TargetMacosArch   `json:"amd64,omitempty"`
}

func (t *TargetMacos) GetMacosArches() map[consts.Arch]*TargetMacosArch {
	return map[consts.Arch]*TargetMacosArch{
		consts.ArchArm64: &t.Arm64,
		consts.ArchAmd64: &t.Amd64,
	}
}

func (t *TargetMacos) GetArches() map[consts.Arch]ArchTarget {
	return asArchTargets(t.GetMacosArches())
}

func (t *TargetMacos) GetCommonCfg() TargetsCommon {
	return t.Common.TargetsCommon
}

func (t *TargetMacos) SetCommonCfg(cfg TargetsCommon) {
	t.Common.TargetsCommon = cfg
}

func (t *TargetMacos) ShouldBuildAnyBundle() bool {
	return maputil.AnyValue(t.GetMacosArches(), func(arch *TargetMacosArch) bool {
		return arch.ShouldBuildBundle()
	})
}
