package bundle

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/consts"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	"github.com/cardinalby/xgo-pack/pkg/platforms/macos/iconset"
	"github.com/cardinalby/xgo-pack/pkg/platforms/macos/plist"
	fsutil "github.com/cardinalby/xgo-pack/pkg/util/fs"
	"github.com/cardinalby/xgo-pack/pkg/util/fs/fs_resource"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
	"golang.org/x/sync/errgroup"
)

const defaultIconSetName = "icon.icns"

func RegisterBuilders(ctx buildctx.Context) {
	ctx.Artifacts.RegisterBuilder(buildctx.KindMacosIconSet, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		pngIcon, err := ctx.Artifacts.Get(ctx, buildctx.KindPngIcon)
		if err != nil {
			return nil, err
		}
		image, err := getPngImage(ctx, pngIcon.GetPath())
		if err != nil {
			return nil, err
		}

		iconSetData, err := iconset.Generate(image)
		if err != nil {
			return nil, err
		}

		iconSetFile, err := ctx.NewTempFile("macos-icon-set.icns")
		if err != nil {
			return nil, fmt.Errorf("error creating macos iconset temp file path: %w", err)
		}

		if err := os.WriteFile(iconSetFile.GetPath(), iconSetData, 0755); err != nil {
			return nil, fmt.Errorf("error writing macos iconset file: %w", err)
		}

		return iconSetFile, nil
	})

	ctx.Artifacts.RegisterBuilder(buildctx.KindMacosPlist, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		plistProps := plist.Properties{
			ProductName:             ctx.Cfg.Targets.Macos.Common.ProductName,
			BinName:                 ctx.Cfg.Targets.Common.BinName,
			Identifier:              ctx.Cfg.Targets.Macos.Common.Identifier,
			ProductVersion:          ctx.Cfg.Targets.Macos.Common.Version,
			IconFile:                defaultIconSetName,
			IsHighResolutionCapable: typeutil.PtrValueOrDefault(ctx.Cfg.Targets.Macos.Common.HighDpi),
			IsHideInDock:            typeutil.PtrValueOrDefault(ctx.Cfg.Targets.Macos.Common.Bundle.HideInDock),
			CopyRight:               ctx.Cfg.Targets.Macos.Common.Copyright,
		}
		plistData, err := plist.GetPlist(plistProps)
		if err != nil {
			return nil, fmt.Errorf("error creating macos plist: %w", err)
		}
		plistFile, err := ctx.NewTempFile("macos-plist.plist")
		if err != nil {
			return nil, fmt.Errorf("error creating macos plist temp file path: %w", err)
		}
		if err := fsutil.WriteFile(plistFile.GetPath(), plistData); err != nil {
			return nil, fmt.Errorf("error writing macos plist file: %w", err)
		}
		return plistFile, nil
	})

	for arch, archCfg := range ctx.Cfg.Targets.Macos.GetMacosArches() {
		arch := arch
		archCfg := archCfg
		outDir := filepath.Join(ctx.Cfg.GetDistDirPath(), archCfg.GetOutDir())
		ctx.Artifacts.RegisterBuilder(buildctx.MacosBundleKind(arch), func(ctx buildctx.Context) (buildctx.Artifact, error) {
			var bundleDir buildctx.Artifact
			if archCfg.ShouldKeepBundle() {
				bundleDir = fs_resource.NewPermanentDir(filepath.Join(
					outDir,
					ctx.Cfg.Targets.Macos.Common.Bundle.BundleName,
				))
			} else {
				var err error
				bundleDir, err = fs_resource.NewTempDir(outDir, ctx.Cfg.Targets.Macos.Common.Bundle.BundleName)
				if err != nil {
					return nil, err
				}
			}
			outBundlePath := Path(bundleDir.GetPath())

			var plistFile, iconSetFile, binFile buildctx.Artifact

			errGr := errgroup.Group{}

			errGr.Go(func() (err error) {
				plistFile, err = ctx.Artifacts.Get(ctx, buildctx.KindMacosPlist)
				return err
			})

			errGr.Go(func() (err error) {
				iconSetFile, err = ctx.Artifacts.Get(ctx, buildctx.KindMacosIconSet)
				return err
			})

			errGr.Go(func() (err error) {
				binFile, err = ctx.Artifacts.Get(ctx, buildctx.BinKind(consts.OsDarwin, arch))
				return err
			})

			if err := errGr.Wait(); err != nil {
				return nil, err
			}

			dstBinPath := filepath.Join(outBundlePath.GetBinDir(), ctx.Cfg.Targets.Common.BinName)
			if err := fsutil.CopyFile(binFile.GetPath(), dstBinPath); err != nil {
				return nil, err
			}

			if ctx.Err() != nil {
				return nil, ctx.Err()
			}

			dstIconSetPath := filepath.Join(outBundlePath.GetResourcesPath(), defaultIconSetName)
			if err := fsutil.CopyFile(iconSetFile.GetPath(), dstIconSetPath); err != nil {
				return nil, err
			}

			if ctx.Err() != nil {
				return nil, ctx.Err()
			}

			if err := fsutil.CopyFile(plistFile.GetPath(), outBundlePath.GetPlistPath()); err != nil {
				return nil, err
			}

			return bundleDir, nil
		})
	}
}
