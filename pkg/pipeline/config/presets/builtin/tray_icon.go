package builtin

import (
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

const TrayIconPresetName cfgtypes.PresetName = "xgo-pack:tray_icon"

func init() {
	Presets[TrayIconPresetName] = trayIconCfg
}

var trayIconCfg = cfgtypes.Config{
	Presets: []cfgtypes.PresetName{GuiPresetName},
	Targets: cfgtypes.Targets{
		Macos: cfgtypes.TargetMacos{
			Common: cfgtypes.TargetMacosCommon{
				Bundle: cfgtypes.TargetMacosCommonBundle{
					HideInDock: typeutil.Ptr(true),
				},
			},
		},
	},
}
