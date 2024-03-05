package cfgtypes

import (
	"github.com/cardinalby/xgo-pack/pkg/build_go/config"
)

type TargetsCommon struct {
	// ProductName is a human-readable name of the product. Is used in MacOS app bundle, dmg and Windows manifest
	// If not set, the last part of module name will be used
	ProductName string `json:"product_name,omitempty"`
	// Version is a version of the app. It's used in MacOS app bundle and Windows manifest
	// If not set, "1.0.0" will be used
	Version string `json:"version,omitempty"`
	// Identifier is a unique identifier for the app (used in MacOS bundle plist and Windows manifest),
	// usually in reverse domain notation e.g. com.example.myapp
	// If not set, the reversed module name will be used
	Identifier string `json:"identifier,omitempty"`
	// Copyright is a copyright string (used in MacOS bundle plist and deb package)
	// If not set, the "© [current_year], [identifier without the last part]" will be used
	Copyright string `json:"copyright,omitempty"`
	// HighDpi is a flag to enable high dpi support on Windows and MacOS
	HighDpi *bool `json:"high_dpi,omitempty"`
	// Arguments of go build command
	GoBuild config.TargetBuildConfig `json:"go_build,omitempty"`
	// The name of the resulting binary file. For Windows '.exe' extension will be added automatically
	// If not set, the last part of the main package path will be used
	BinName string `json:"bin_name,omitempty"`
}
