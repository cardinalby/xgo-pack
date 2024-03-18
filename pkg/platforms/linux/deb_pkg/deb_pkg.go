package deb_pkg

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/consts"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	"github.com/cardinalby/xgo-pack/pkg/util/fs/fs_resource"
	"github.com/goccy/go-yaml"
	"github.com/goreleaser/nfpm/v2"
	_ "github.com/goreleaser/nfpm/v2/deb"
	"github.com/goreleaser/nfpm/v2/files"
)

func RegisterDebPkgBuilders(ctx buildctx.Context) {
	registerDesktopEntryBuilder(ctx)

	for arch, archCfg := range ctx.Cfg.Targets.Linux.GetArches() {
		arch := arch
		archCfg := archCfg
		ctx.Artifacts.RegisterBuilder(buildctx.LinuxDebKind(arch), func(ctx buildctx.Context) (buildctx.Artifact, error) {
			outPath := filepath.Join(
				ctx.Cfg.GetDistDirPath(), archCfg.GetOutDir(), ctx.Cfg.Targets.Linux.Common.Deb.DebName,
			)

			bin, err := ctx.Artifacts.Get(ctx, buildctx.BinKind(consts.OsLinux, arch))
			if err != nil {
				return nil, fmt.Errorf("error getting %s/%s bin: %w", consts.OsLinux, arch, err)
			}

			srcIcon, err := ctx.Artifacts.Get(ctx, buildctx.KindPngIcon)
			if err != nil {
				return nil, fmt.Errorf("error getting png icon: %w", err)
			}

			srcDesktopEntry, err := ctx.Artifacts.Get(ctx, buildctx.KindLinuxDesktopEntry)
			if err != nil {
				return nil, fmt.Errorf("error getting desktop entry: %w", err)
			}

			if err := BuildDebPkg(
				ctx,
				arch,
				bin.GetPath(),
				srcIcon.GetPath(),
				srcDesktopEntry.GetPath(),
				outPath,
			); err != nil {
				return nil, fmt.Errorf("error building deb pkg: %w", err)
			}

			return fs_resource.NewPermanentFile(outPath), nil
		})
	}

}

func BuildDebPkg(
	ctx buildctx.Context,
	arch consts.Arch,
	srcBinPath string,
	srcIconPath string,
	srcDesktopEntryPath string,
	outPath string,
) error {
	packager, err := nfpm.Get("deb")
	if err != nil {
		return fmt.Errorf("error getting deb packager: %w", err)
	}

	nfpmCfg, err := getNfpmConfig(ctx, arch, srcBinPath, srcIconPath, srcDesktopEntryPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("error creating deb pkg dir: %w", err)
	}
	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("error creating deb pkg file '%s': %w", outPath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			ctx.Logger.Printf("error closing deb pkg file '%s': %v", outPath, err)
		}
	}()

	if err := packager.Package(&nfpmCfg.Info, file); err != nil {
		return fmt.Errorf("error packaging deb: %w", err)
	}

	return nil
}

func getNfpmConfig(
	ctx buildctx.Context,
	arch consts.Arch,
	srcBinPath string,
	srcIconPath string,
	srcDesktopEntryPath string,
) (cfg nfpm.Config, err error) {
	debCfg := ctx.Cfg.Targets.Linux.Common.Deb
	envMapping, err := getCfgEnvMapping(ctx, arch)
	if err != nil {
		return cfg, err
	}
	if debCfg.CustomNfpmConfig != "" {
		customCfgAbsPath := filepath.Join(ctx.Cfg.Root, debCfg.CustomNfpmConfig)
		if cfg, err = nfpm.ParseFileWithEnvMapping(customCfgAbsPath, envMapping); err != nil {
			return cfg, fmt.Errorf(
				"error parsing custom '%s' nfpm config file: %w",
				customCfgAbsPath, err,
			)
		}
		return cfg, nil
	}

	contents := files.Contents{
		{
			Source:      srcBinPath,
			Destination: debCfg.DstBinPath,
		},
		{
			Source: srcDesktopEntryPath,
			Destination: filepath.Join(
				"/usr/share/applications/",
				ctx.Cfg.Targets.Common.ProductName+".desktop",
			),
		},
	}
	if srcIconPath != "" {
		contents = append(contents, &files.Content{
			Source:      srcIconPath,
			Destination: debCfg.DesktopEntry.DstIconPath,
		})
	}

	for src, dst := range debCfg.Contents {
		if err != nil {
			return cfg, fmt.Errorf("error joining root and src path: %w", err)
		}
		contents = append(contents, &files.Content{
			Source:      src,
			Destination: dst,
		})
	}

	baseCfg := nfpm.Config{
		Info: nfpm.Info{
			Name:        "${XGO_PACK_DEB_NAME}",
			Arch:        "${XGO_PACK_ARCH}",
			Version:     "${XGO_PACK_VERSION}",
			Section:     debCfg.Section,
			Maintainer:  debCfg.Maintainer,
			Description: debCfg.Description,
			Vendor:      debCfg.Vendor,
			Homepage:    debCfg.Homepage,
			License:     debCfg.License,
			Overridables: nfpm.Overridables{
				Contents: contents,
			},
		},
	}
	baseCfgYamlData, err := yaml.Marshal(baseCfg)
	if err != nil {
		return cfg, fmt.Errorf("error marshalling base nfpm config: %w", err)
	}
	buff := bytes.NewBuffer(baseCfgYamlData)
	return nfpm.ParseWithEnvMapping(buff, envMapping)
}

func getCfgEnvMapping(
	ctx buildctx.Context,
	arch consts.Arch,
) (func(string) string, error) {
	debName := toDebName(ctx.Cfg.Targets.Linux.Common.Deb.Name)
	if err := validateDebName(debName); err != nil {
		return nil, err
	}

	m := map[string]string{
		"XGO_PACK_DEB_NAME": debName,
		"XGO_PACK_ARCH":     string(arch),
		"XGO_PACK_VERSION":  ctx.Cfg.Targets.Linux.Common.Version,
	}
	return func(key string) string {
		return m[key]
	}, nil
}
