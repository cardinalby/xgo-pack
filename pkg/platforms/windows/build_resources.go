package windows

import (
	"fmt"
	"path/filepath"

	"github.com/akavel/rsrc/rsrc"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	"github.com/cardinalby/xgo-pack/pkg/platforms/windows/manifest"
	"github.com/cardinalby/xgo-pack/pkg/util/fs/fs_resource"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
	"golang.org/x/sync/errgroup"
)

func RegisterBuilders(ctx buildctx.Context) {
	ctx.Artifacts.RegisterBuilder(buildctx.KindWinManifest, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		manifestProperties := manifest.Properties{
			Identifier: ctx.Cfg.Targets.Windows.Common.Identifier,
			IsDpiAware: typeutil.PtrValueOrDefault(ctx.Cfg.Targets.Windows.Common.HighDpi),
		}

		manifestData, err := manifest.GetManifest(manifestProperties)
		if err != nil {
			return nil, fmt.Errorf("error preparing manifest: %w", err)
		}

		manifestTmpFile, err := ctx.WriteTempFile(manifestData, "manifest.xml")
		if err != nil {
			return nil, fmt.Errorf("error saving temp manifest file: %w", err)
		}

		return manifestTmpFile, nil
	})

	for arch, archCfg := range ctx.Cfg.Targets.Windows.GetWinArches() {
		arch := arch
		archCfg := archCfg
		ctx.Artifacts.RegisterBuilder(buildctx.WinSysoKind(arch), func(ctx buildctx.Context) (buildctx.Artifact, error) {
			var sysoFile buildctx.Artifact
			sysoFilePath := filepath.Join(ctx.Cfg.Root, ctx.Cfg.Src.MainPkg, fmt.Sprintf("resources_%s.syso", arch))
			// Don't create a new syso file because multiple syso files cause "too many .rsrc sections" build error
			if archCfg.ShouldKeepSyso() {
				sysoFile = fs_resource.NewPermanentFile(sysoFilePath)
			} else {
				sysoFile = &fs_resource.TempFile{Path: sysoFilePath}
			}

			errGr := &errgroup.Group{}

			var icoFile, manifestFile buildctx.Artifact

			errGr.Go(func() (err error) {
				icoFile, err = ctx.Artifacts.Get(ctx, buildctx.KindIcoIcon)
				return
			})

			errGr.Go(func() (err error) {
				manifestFile, err = ctx.Artifacts.Get(ctx, buildctx.KindWinManifest)
				return
			})

			if err := errGr.Wait(); err != nil {
				return nil, err
			}

			if err := rsrc.Embed(sysoFilePath, string(arch), manifestFile.GetPath(), icoFile.GetPath()); err != nil {
				return nil, fmt.Errorf("error creating '%s' file: %w", sysoFilePath, err)
			}

			return sysoFile, nil
		})
	}
}
