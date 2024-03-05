package builders

import (
	"fmt"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/imagemagick"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
)

func RegisterBuilders(ctx buildctx.Context) {
	ctx.Artifacts.RegisterBuilder(buildctx.KindPngIcon, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		iconArtifact, err := ctx.Artifacts.Get(ctx, buildctx.KindIcon)
		if err != nil {
			return nil, fmt.Errorf("error getting icon artifact: %w", err)
		}
		if filepath.Ext(iconArtifact.GetPath()) == ".png" {
			return iconArtifact, nil
		}
		tempPngIcon, err := ctx.NewTempFile("icon.png")
		if err != nil {
			return nil, fmt.Errorf("error creating temp png icon file path: %w", err)
		}
		if err := imagemagick.ToPng(ctx, buildctx.ImgSourcePath(iconArtifact.GetPath()), tempPngIcon.Path); err != nil {
			return nil, fmt.Errorf("error converting icon to png: %w", err)
		}
		return tempPngIcon, nil
	})

	ctx.Artifacts.RegisterBuilder(buildctx.KindIcoIcon, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		iconArtifact, err := ctx.Artifacts.Get(ctx, buildctx.KindIcon)
		if err != nil {
			return nil, fmt.Errorf("error getting icon artifact: %w", err)
		}
		if filepath.Ext(iconArtifact.GetPath()) == ".ico" {
			return iconArtifact, nil
		}
		tempIcoIcon, err := ctx.NewTempFile("icon.ico")
		if err != nil {
			return nil, fmt.Errorf("error creating temp ico icon file path: %w", err)
		}
		if err := imagemagick.ToIco(
			ctx,
			buildctx.ImgSourcePath(iconArtifact.GetPath()),
			tempIcoIcon.Path,
			[]int{16, 32, 48, 64, 128, 256},
		); err != nil {
			return nil, fmt.Errorf("error converting icon to ico: %w", err)
		}
		return tempIcoIcon, nil
	})
}
