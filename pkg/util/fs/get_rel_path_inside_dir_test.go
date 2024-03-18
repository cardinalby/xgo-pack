package fsutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRelPathInsideDir(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		basePath      string
		targetPath    string
		want          string
		wantNotRelErr bool
		wantErr       bool
	}{
		{
			name:          "base path is prefix of target path",
			basePath:      "/a/b",
			targetPath:    "/a/b/c",
			want:          "c",
			wantNotRelErr: false,
			wantErr:       false,
		},
		{
			name:          "base path is not prefix of target path",
			basePath:      "/a/b",
			targetPath:    "/a/c",
			want:          "",
			wantNotRelErr: true,
			wantErr:       true,
		},
		{
			name:          "base path is prefix of target path",
			basePath:      "/a/b",
			targetPath:    "/a/b/c/d/..",
			want:          "c",
			wantNotRelErr: false,
			wantErr:       false,
		},
		{
			name:          "base path is not prefix of target path",
			basePath:      "/a/b",
			targetPath:    "/a/b/../",
			want:          "",
			wantNotRelErr: true,
			wantErr:       true,
		},
		{
			name:          "different roots",
			basePath:      "/a/b",
			targetPath:    "/c/d",
			want:          "",
			wantNotRelErr: true,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetRelPathInsideDir(tt.targetPath, tt.basePath)
			if tt.wantErr {
				require.Error(t, err)
			}
			if tt.wantNotRelErr {
				require.ErrorIs(t, err, ErrIsNotRelative)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
