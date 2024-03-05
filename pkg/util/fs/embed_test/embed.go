package embed_test

import "embed"

//go:embed a.txt
//go:embed nested
var EmbedFS embed.FS
