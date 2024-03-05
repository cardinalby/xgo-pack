package util

import (
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
)

func ReadConfigFile(filePath string) (cfgtypes.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfgtypes.Config{}, fmt.Errorf("error reading '%s': %w", filePath, err)
	}
	var presetCfg cfgtypes.Config
	if err := yaml.Unmarshal(data, &presetCfg); err != nil {
		return cfgtypes.Config{}, fmt.Errorf("error parsing file '%s': %w", filePath, err)
	}
	return presetCfg, nil
}
