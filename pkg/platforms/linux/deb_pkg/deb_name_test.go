package deb_pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateDebName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		debName string
		wantErr bool
	}{
		{
			name:    "valid",
			debName: "my-package",
			wantErr: false,
		},
		{
			name:    "invalid",
			debName: "MyPackage",
			wantErr: true,
		},
		{
			name:    "invalid",
			debName: "my_package",
			wantErr: true,
		},
		{
			name:    "invalid",
			debName: "my package",
			wantErr: true,
		},
		{
			name:    "invalid",
			debName: "1ab",
		},
		{
			name:    "invalid",
			debName: "a",
			wantErr: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := validateDebName(test.debName)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToDebName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no change",
			in:   "my-package",
			want: "my-package",
		},
		{
			name: "replace space",
			in:   "my package",
			want: "my-package",
		},
		{
			name: "replace invalid char",
			in:   "my_package",
			want: "my-package",
		},
		{
			name: "replace invalid char",
			in:   "MyPackage",
			want: "mypackage",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.want, toDebName(test.in))
		})
	}
}
