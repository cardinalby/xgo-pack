package pipeline

import (
	"context"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/consts"
	"github.com/cardinalby/xgo-pack/pkg/platforms/linux/deb_pkg"
	"github.com/cardinalby/xgo-pack/pkg/util/fs/fs_resource"
	"golang.org/x/sync/errgroup"

	go_builders "github.com/cardinalby/xgo-pack/pkg/build_go/builders"
	imagemagick_builders "github.com/cardinalby/xgo-pack/pkg/imagemagick/builders"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	"github.com/cardinalby/xgo-pack/pkg/platforms/macos/bundle"
	dmg_builders "github.com/cardinalby/xgo-pack/pkg/platforms/macos/dmg/builders"
	"github.com/cardinalby/xgo-pack/pkg/platforms/windows"
	"github.com/cardinalby/xgo-pack/pkg/resource/default_icon"
	"github.com/cardinalby/xgo-pack/pkg/util/logging"
)

func Start(
	ctx context.Context,
	logger logging.Logger,
	cfg cfgtypes.Config,
) (err error) {
	logger = logging.NewSyncedWrapper(logger)
	errGroup, ctx := errgroup.WithContext(ctx)
	buildCtx := buildctx.NewContext(ctx, cfg, buildctx.NewArtifacts(logger), logger)
	defer func() {
		if err := buildCtx.Artifacts.Dispose(); err != nil {
			logger.Printf("Failed to dispose artifacts: %v", err)
		}
	}()
	registerArtifactBuilders(buildCtx)

	for os, osTargets := range cfg.Targets.GetOsArches() {
		os := os
		osTargets := osTargets
		for arch, archCfg := range osTargets.GetArches() {
			arch := arch
			archCfg := archCfg
			if archCfg.ShouldBuildBin() {
				errGroup.Go(func() error {
					if os == consts.OsWindows {
						if _, err := buildCtx.Artifacts.Get(buildCtx, buildctx.WinSysoKind(arch)); err != nil {
							return err
						}
					}
					_, err := buildCtx.Artifacts.Get(buildCtx, buildctx.BinKind(os, arch))
					return err
				})
			}
			if macosArchCfg, ok := archCfg.(*cfgtypes.TargetMacosArch); ok {
				if macosArchCfg.ShouldBuildBundle() {
					errGroup.Go(func() error {
						_, err := buildCtx.Artifacts.Get(buildCtx, buildctx.MacosBundleKind(arch))
						return err
					})
				}

				if macosArchCfg.ShouldBuildDmg() {
					errGroup.Go(func() error {
						_, err := buildCtx.Artifacts.Get(buildCtx, buildctx.MacosDmgKind(arch))
						return err
					})
				}
			}
			if linuxArchCfg, ok := archCfg.(*cfgtypes.TargetLinuxArch); ok {
				if linuxArchCfg.ShouldBuildDeb() {
					errGroup.Go(func() error {
						_, err := buildCtx.Artifacts.Get(buildCtx, buildctx.LinuxDebKind(arch))
						return err
					})
				}
			}
		}
	}

	return errGroup.Wait()
}

func registerArtifactBuilders(ctx buildctx.Context) {
	default_icon.RegisterBuilder(ctx)
	ctx.Artifacts.RegisterBuilder(buildctx.KindIcon, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		if ctx.Cfg.Src.Icon == cfgtypes.DefaultBuiltInIcon {
			return ctx.Artifacts.Get(ctx, buildctx.KindDefaultPngIcon)
		}
		return fs_resource.NewPermanentFile(filepath.Join(ctx.Cfg.Root, ctx.Cfg.Src.Icon)), nil
	})
	imagemagick_builders.RegisterBuilders(ctx)
	windows.RegisterBuilders(ctx)
	go_builders.RegisterBuilders(ctx)
	bundle.RegisterBuilders(ctx)
	dmg_builders.RegisterBuilders(ctx)
	deb_pkg.RegisterDebPkgBuilders(ctx)
}
