package builders

import (
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/build_go"
	"github.com/cardinalby/xgo-pack/pkg/build_go/config"
	"github.com/cardinalby/xgo-pack/pkg/consts"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	"github.com/cardinalby/xgo-pack/pkg/util/fs/fs_resource"
)

func RegisterBuilders(ctx buildctx.Context) {
	ctx.Artifacts.RegisterBuilder(buildctx.KindBinTempDir, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		binTmpDir, err := ctx.NewTempDir("xgo_bin")
		if err != nil {
			return nil, err
		}
		return &fs_resource.TempDir{Path: binTmpDir.GetPath()}, nil
	})

	for os, osTargets := range ctx.Cfg.Targets.GetOsArches() {
		os := os
		osCommonCfg := ctx.Cfg.Targets.GetOsCommonCfg(os)
		binFileName := ctx.Cfg.Targets.Common.BinName
		if os == consts.OsWindows {
			binFileName += ".exe"
		}

		for arch, archCfg := range osTargets.GetArches() {
			arch := arch
			archCfg := archCfg
			ctx.Artifacts.RegisterBuilder(
				buildctx.BinKind(os, arch),
				func(ctx buildctx.Context) (buildctx.Artifact, error) {
					outDir := filepath.Join(
						ctx.Cfg.GetDistDirPath(),
						archCfg.GetOutDir(),
					)
					binFilePath := filepath.Join(outDir, binFileName)
					var binFile buildctx.Artifact
					if archCfg.ShouldKeepBin() {
						binFile = fs_resource.NewPermanentFile(binFilePath)
					} else {
						binFile = &fs_resource.TempFile{Path: binFilePath}
					}

					binTempDir, err := ctx.Artifacts.Get(ctx, buildctx.KindBinTempDir)
					if err != nil {
						return nil, err
					}

					buildCfg := config.Config{
						XGoConfig:      ctx.Cfg.XGo,
						RootPath:       ctx.Cfg.Root,
						BinTempDir:     binTempDir.GetPath(),
						MainPkgRelPath: ctx.Cfg.Src.MainPkg,
						Targets: map[config.Target]config.TargetConfig{
							config.Target{
								Os:   os,
								Arch: arch,
							}: {
								OutBinPath: binFilePath,
								TargetBuildConfig: config.TargetBuildConfig{
									Race:      osCommonCfg.GoBuild.Race,
									Tags:      osCommonCfg.GoBuild.Tags,
									LdFlags:   osCommonCfg.GoBuild.LdFlags,
									Mode:      osCommonCfg.GoBuild.Mode,
									VCS:       osCommonCfg.GoBuild.VCS,
									TrimPath:  osCommonCfg.GoBuild.TrimPath,
									CrossArgs: osCommonCfg.GoBuild.CrossArgs,
								},
							},
						},
					}

					if err := build_go.Start(ctx, buildCfg); err != nil {
						return nil, err
					}

					return binFile, nil
				})
		}
	}
}
