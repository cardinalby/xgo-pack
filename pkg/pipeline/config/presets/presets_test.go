package presets

import (
	"testing"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	"github.com/stretchr/testify/require"
)

func TestApplyConfigPreset(t *testing.T) {
	p1 := cfgtypes.Config{
		Root: "p1",
	}
	p2 := cfgtypes.Config{
		Root: "p2",
	}
	p12 := cfgtypes.Config{
		Presets: []cfgtypes.PresetName{"p1", "p2"},
	}
	p21 := cfgtypes.Config{
		Presets: []cfgtypes.PresetName{"p2", "p1"},
	}
	p3 := cfgtypes.Config{
		Presets: []cfgtypes.PresetName{"p21"},
	}
	testResolver := Map{
		"p1":  p1,
		"p2":  p2,
		"p12": p12,
		"p21": p21,
		"p3":  p3,
	}

	res, err := ApplyConfigPreset(cfgtypes.Config{
		Presets: []cfgtypes.PresetName{"p3"},
	}, testResolver)
	require.NoError(t, err)
	require.Equal(t, cfgtypes.Config{
		Root: "p1",
	}, res)

	res, err = ApplyConfigPreset(cfgtypes.Config{
		Presets: []cfgtypes.PresetName{"p12"},
	}, testResolver)
	require.NoError(t, err)
	require.Equal(t, cfgtypes.Config{
		Root: "p2",
	}, res)

	res, err = ApplyConfigPreset(cfgtypes.Config{
		Presets: []cfgtypes.PresetName{"x"},
	}, testResolver)
	require.Error(t, err)
}
