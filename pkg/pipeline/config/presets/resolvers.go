package presets

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/util"
)

type Map map[cfgtypes.PresetName]cfgtypes.Config

func (p Map) GetPresetConfig(preset cfgtypes.PresetName) (cfgtypes.Config, error) {
	if cfg, ok := p[preset]; ok {
		return cfg, nil
	}
	return cfgtypes.Config{}, fmt.Errorf("unknown preset: '%s'", preset)
}

func (p Map) GetNames() []cfgtypes.PresetName {
	var names []cfgtypes.PresetName
	for name := range p {
		names = append(names, name)
	}
	return names
}

type Dir string

func (d Dir) GetPresetConfig(preset cfgtypes.PresetName) (cfgtypes.Config, error) {
	filePath := filepath.Join(string(d), string(preset))
	return util.ReadConfigFile(filePath)
}

type Compound []Resolver

func (c Compound) GetPresetConfig(preset cfgtypes.PresetName) (cfgtypes.Config, error) {
	var errs []error
	for _, r := range c {
		cfg, err := r.GetPresetConfig(preset)
		if err == nil {
			return cfg, nil
		}
		errs = append(errs, err)
	}
	return cfgtypes.Config{}, errors.Join(errs...)
}
