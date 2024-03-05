package builtin

import (
	"testing"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/presets"
	"github.com/stretchr/testify/require"
)

func TestPresets(t *testing.T) {
	t.Parallel()
	for _, p := range Presets {
		_, err := presets.ApplyConfigPreset(p, Presets)
		require.NoError(t, err)
	}
}
