package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModulePathToIdentifier(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			"simple",
			"github.com/cardinalby/xgo-pack",
			"com.github.cardinalby.xgo-pack",
		},
		{
			"with subpackage",
			"github.com/cardinalby/xgo-pack/pkg/util",
			"com.github.cardinalby.xgo-pack.pkg.util",
		},
		{
			"with no domain",
			"custom",
			"custom",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, ModulePathToIdentifier(tt.path))
		})
	}
}

func TestIdentifierLastPart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			"simple",
			"com.github.cardinalby.xgo-pack",
			"xgo-pack",
		},
		{
			"with subpackage",
			"com.github.cardinalby.xgo-pack.pkg.util",
			"util",
		},
		{
			"with no domain",
			"custom",
			"custom",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, IdentifierLastPart(tt.in))
		})
	}
}

func TestIdentifierWithoutLastPart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			"simple",
			"com.github.cardinalby.xgo-pack",
			"com.github.cardinalby",
		},
		{
			"with subpackage",
			"com.github.cardinalby.xgo-pack.pkg.util",
			"com.github.cardinalby.xgo-pack.pkg",
		},
		{
			"with no domain",
			"custom",
			"custom",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, IdentifierWithoutLastPart(tt.in))
		})
	}
}
