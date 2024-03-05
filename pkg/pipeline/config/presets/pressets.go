package presets

import (
	"errors"
	"fmt"

	dtomerge "github.com/cardinalby/go-dto-merge"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
)

type Resolver interface {
	GetPresetConfig(preset cfgtypes.PresetName) (cfgtypes.Config, error)
}

func ApplyConfigPreset(
	cfg cfgtypes.Config,
	resolver Resolver,
) (cfgtypes.Config, error) {
	return applyConfigPresetImpl(cfg, resolver, make(map[cfgtypes.PresetName]struct{}))
}

func applyConfigPresetImpl(
	cfg cfgtypes.Config,
	resolver Resolver,
	visitedPresets map[cfgtypes.PresetName]struct{},
) (cfgtypes.Config, error) {
	presets := cfg.Presets
	if len(presets) == 0 {
		return cfg, nil
	}

	if resolver == nil {
		return cfgtypes.Config{}, errors.New("presets resolver is not set")
	}

	presetConfigs := make([]cfgtypes.Config, 0, len(presets))
	for _, preset := range presets {
		if _, has := visitedPresets[preset]; has {
			return cfgtypes.Config{}, fmt.Errorf("loop in presets: preset '%s' is already visited", preset)
		}
		visitedPresets[preset] = struct{}{}
		presetCfg, err := resolver.GetPresetConfig(preset)
		if err != nil {
			return cfgtypes.Config{}, fmt.Errorf("error resolving preset '%s': %w", preset, err)
		}
		presetCfg, err = applyConfigPresetImpl(presetCfg, resolver, visitedPresets)
		if err != nil {
			return cfgtypes.Config{}, fmt.Errorf("error processing preset '%s': %w", preset, err)
		}
		presetConfigs = append(presetConfigs, presetCfg)
	}

	var err error
	resultingPresetCfg := presetConfigs[0]
	for i := 1; i < len(presetConfigs); i++ {
		presetCfg := presetConfigs[i]
		resultingPresetCfg, err = dtomerge.Merge(resultingPresetCfg, presetCfg)
		if err != nil {
			return cfgtypes.Config{}, fmt.Errorf("error applying preset '%s': %w", presets, err)
		}
	}

	cfg, err = dtomerge.Merge(resultingPresetCfg, cfg)
	if err != nil {
		return cfgtypes.Config{}, fmt.Errorf("error applying config on top of '%v' presets: %w", presets, err)
	}
	cfg.Presets = nil
	return cfg, nil
}
