package config

import (
	"fmt"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/presets"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/presets/builtin"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/util"
)

func ReadFile(filePath string, presetsDir string) (cfg cfgtypes.Config, err error) {
	presetsResolver := getDefaultCfgPresetsResolver(presetsDir)
	cfg, err = util.ReadConfigFile(filePath)
	if err != nil {
		return cfg, err
	}
	if cfg, err = Prepare(cfg, presetsResolver); err != nil {
		return cfgtypes.Config{}, fmt.Errorf("error preparing config: %w", err)
	}
	return cfg, nil
}

func getDefaultCfgPresetsResolver(presetsDir string) presets.Resolver {
	if presetsDir == "" {
		return builtin.Presets
	}
	return presets.Compound{
		builtin.Presets,
		presets.Dir(presetsDir),
	}
}
