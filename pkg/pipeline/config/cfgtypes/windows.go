package cfgtypes

import (
	"github.com/cardinalby/xgo-pack/pkg/consts"
	"github.com/cardinalby/xgo-pack/pkg/util/maputil"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

type TargetWindowsArch struct {
	// OutDir is a path relative to Config.DistDir where the final artifacts will be placed
	// If not set, "windows_[arch]" will be used
	OutDir string `json:"out_dir,omitempty"`
	// BuildSyso defines if .syso file should be built.
	// If false, the .syso file will be built only in case it's required to build binary
	// Temporary path will be used in this case.
	// To keep the .syso file, set it to true
	BuildSyso *bool `json:"build_syso,omitempty"`
	// Defines if binary should be built
	BuildBin *bool `json:"build_bin,omitempty"`
}

func (t *TargetWindowsArch) ShouldBuildSyso() bool {
	return typeutil.PtrValueOrDefault(t.BuildSyso) || typeutil.PtrValueOrDefault(t.BuildBin)
}

func (t *TargetWindowsArch) ShouldKeepSyso() bool {
	return typeutil.PtrValueOrDefault(t.BuildSyso)
}

func (t *TargetWindowsArch) ShouldBuildBin() bool {
	return typeutil.PtrValueOrDefault(t.BuildBin)
}

func (t *TargetWindowsArch) ShouldKeepBin() bool {
	return typeutil.PtrValueOrDefault(t.BuildBin)
}

func (t *TargetWindowsArch) GetOutDir() string {
	return t.OutDir
}

func (t *TargetWindowsArch) SetOutDir(outDir string) {
	t.OutDir = outDir
}

type TargetWindowsCommon struct {
	TargetsCommon `json:",inline"`
}

type TargetWindows struct {
	// Common will be used as defaults for all architectures
	Common TargetWindowsCommon `json:"common,omitempty"`
	Amd64  TargetWindowsArch   `json:"amd64,omitempty"`
}

func (t *TargetWindows) GetWinArches() map[consts.Arch]*TargetWindowsArch {
	return map[consts.Arch]*TargetWindowsArch{
		consts.ArchAmd64: &t.Amd64,
	}
}

func (t *TargetWindows) GetArches() map[consts.Arch]ArchTarget {
	return asArchTargets(t.GetWinArches())
}

func (t *TargetWindows) GetCommonCfg() TargetsCommon {
	return t.Common.TargetsCommon
}

func (t *TargetWindows) SetCommonCfg(cfg TargetsCommon) {
	t.Common.TargetsCommon = cfg
}

func (t *TargetWindows) ShouldBuildAnySyso() bool {
	return maputil.AnyValue(t.GetWinArches(), func(arch *TargetWindowsArch) bool {
		return arch.ShouldBuildSyso()
	})
}
