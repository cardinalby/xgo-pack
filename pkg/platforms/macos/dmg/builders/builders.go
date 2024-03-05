package builders

import (
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	"github.com/cardinalby/xgo-pack/pkg/platforms/macos/dmg"
	"github.com/cardinalby/xgo-pack/pkg/platforms/macos/dmg/docker_img"
	"github.com/cardinalby/xgo-pack/pkg/util/fs/fs_resource"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

func RegisterBuilders(ctx buildctx.Context) {

	ctx.Artifacts.RegisterBuilder(
		buildctx.KindMacosCreateDmgDockerImage,
		func(ctx buildctx.Context) (buildctx.Artifact, error) {
			return docker_img.BuildCreateDmgDockerImage(ctx)
		},
	)

	for arch, archCfg := range ctx.Cfg.Targets.Macos.GetArches() {
		arch := arch
		archCfg := archCfg
		ctx.Artifacts.RegisterBuilder(buildctx.MacosDmgKind(arch), func(ctx buildctx.Context) (buildctx.Artifact, error) {
			bundle, err := ctx.Artifacts.Get(ctx, buildctx.MacosBundleKind(arch))
			if err != nil {
				return nil, err
			}

			dmgPath := filepath.Join(
				ctx.Cfg.GetDistDirPath(), archCfg.GetOutDir(), ctx.Cfg.Targets.Macos.Common.Dmg.DmgName,
			)

			if err := dmg.Generate(
				ctx,
				bundle.GetPath(),
				dmgPath,
				ctx.Cfg.Targets.Macos.Common.ProductName,
				typeutil.PtrValueOrDefault(ctx.Cfg.Targets.Macos.Common.Dmg.AddApplicationsSymlink),
			); err != nil {
				return nil, err
			}

			return fs_resource.NewPermanentFile(dmgPath), nil
		})
	}
}
