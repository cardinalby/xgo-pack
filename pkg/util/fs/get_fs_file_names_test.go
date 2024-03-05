package fsutil

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cardinalby/xgo-pack/pkg/util/fs/embed_test"
)

func TestGetFsFilenames(t *testing.T) {
	t.Parallel()
	files, err := GetFsFilenames(embed_test.EmbedFS)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"a.txt", "nested/b.txt"}, files)
}
