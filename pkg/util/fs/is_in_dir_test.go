package fsutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsInDir(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		path string
		dir  string
		want bool
	}{
		{
			name: "nested 2 levels",
			path: "a/b/c",
			dir:  "a",
			want: true,
		},
		{
			name: "nested 1 level",
			path: "a/b/c",
			dir:  "a/b",
			want: true,
		},
		{
			name: "same",
			path: "a/b/c",
			dir:  "a/b/c",
			want: true,
		},
		{
			path: "a/b/c",
			dir:  "a/b/c/d",
			want: true,
		},
		{
			path: "/a",
			dir:  "/",
			want: true,
		},
		{
			path: "/a",
			dir:  "/a",
			want: true,
		},
		{
			path: "/a",
			dir:  "/ab",
			want: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			got, err := IsInDir(tc.dir, tc.path)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
