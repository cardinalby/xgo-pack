package builtin

import (
	"github.com/cardinalby/xgo-pack/pkg/build_go/config"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

const GuiPresetName cfgtypes.PresetName = "xgo-pack:gui"

func init() {
	Presets[GuiPresetName] = guiCfg
}

var guiCfg = cfgtypes.Config{
	Targets: cfgtypes.Targets{
		Common: cfgtypes.TargetsCommon{
			GoBuild: config.TargetBuildConfig{
				LdFlags: "-s -w",
			},
			HighDpi: typeutil.Ptr(true),
		},
		Windows: cfgtypes.TargetWindows{
			Common: cfgtypes.TargetWindowsCommon{
				TargetsCommon: cfgtypes.TargetsCommon{
					GoBuild: config.TargetBuildConfig{
						LdFlags: "-s -w -H windowsgui",
					},
				},
			},
		},
		Macos: cfgtypes.TargetMacos{
			Common: cfgtypes.TargetMacosCommon{
				Dmg: cfgtypes.TargetMacosCommonDmg{
					AddApplicationsSymlink: typeutil.Ptr(true),
				},
			},
		},
		Linux: cfgtypes.TargetLinux{
			Common: cfgtypes.TargetLinuxCommon{
				Deb: cfgtypes.TargetLinuxDeb{
					DesktopEntry: cfgtypes.TargetLinuxDebDesktopEntry{
						Terminal:        typeutil.Ptr(false),
						NoDisplay:       typeutil.Ptr(false),
						AddIcon:         typeutil.Ptr(true),
						AddDesktopEntry: typeutil.Ptr(true),
					},
				},
			},
		},
	},
}
