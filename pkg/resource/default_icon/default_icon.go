package default_icon

import (
	"embed"
	"fmt"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	fsutil "github.com/cardinalby/xgo-pack/pkg/util/fs"
)

//go:embed default_icon.png
var EmbedDefaultIcon embed.FS

const DefaultIconPath = "default_icon.png"

func RegisterBuilder(ctx buildctx.Context) {
	ctx.Artifacts.RegisterBuilder(buildctx.KindDefaultPngIcon, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		tempDefaultIcon, err := ctx.NewTempFile(DefaultIconPath)
		err = fsutil.CopyFsFile(EmbedDefaultIcon, DefaultIconPath, tempDefaultIcon.Path)
		if err != nil {
			return nil, fmt.Errorf("error copying default icon: %w", err)
		}
		return tempDefaultIcon, nil
	})
}
