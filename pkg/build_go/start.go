package build_go

import (
	"fmt"
	"path"
	"path/filepath"

	xgolib "github.com/cardinalby/xgo-as-library"
	"github.com/cardinalby/xgo-pack/pkg/build_go/config"
	"github.com/cardinalby/xgo-pack/pkg/consts"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	fsutil "github.com/cardinalby/xgo-pack/pkg/util/fs"
	"github.com/cardinalby/xgo-pack/pkg/util/logging"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

func Start(
	ctx buildctx.Context,
	config config.Config,
) error {
	outBinPrefix := filepath.Base(config.MainPkgRelPath)

	for target, tConfig := range config.Targets {
		buildArgs := xgolib.BuildArgs{
			Verbose:  typeutil.PtrValueOrDefault(config.XGoConfig.Verbose),
			Race:     typeutil.PtrValueOrDefault(tConfig.TargetBuildConfig.Race),
			Tags:     tConfig.TargetBuildConfig.Tags,
			LdFlags:  tConfig.TargetBuildConfig.LdFlags,
			Mode:     tConfig.TargetBuildConfig.Mode,
			VCS:      tConfig.TargetBuildConfig.VCS,
			TrimPath: typeutil.PtrValueOrDefault(tConfig.TargetBuildConfig.TrimPath),
		}
		args := xgolib.Args{
			Repository: config.RootPath,
			SrcPackage: config.MainPkgRelPath,
			OutFolder:  config.BinTempDir,
			OutPrefix:  outBinPrefix,
			Build:      buildArgs,
			Targets:    []string{target.String()},
		}
		//ctx.Artifacts.AddAnonymous(&resource.TempDir{Path: binTmpDir})
		logger := logging.NewBufferedLogger(ctx.Logger)
		err := xgolib.StartBuild(args, logger)
		logger.Flush()
		if err != nil {
			return fmt.Errorf("error building '%s' target: %w", target.String(), err)
		}

		outBinPath := getXgoOutBinPath(config.BinTempDir, outBinPrefix, target)
		if err := fsutil.RenameFile(outBinPath, tConfig.OutBinPath); err != nil {
			return fmt.Errorf(
				"error moving xgo build result bin '%s' to '%s': %w. It may be a result of unsuccessful build",
				outBinPath, tConfig.OutBinPath, err,
			)
		}
	}

	return nil
}

func getXgoOutBinPath(
	binTmpDir string,
	outPrefix string,
	target config.Target,
) string {
	res := path.Join(binTmpDir, fmt.Sprintf("%s-%s-%s", outPrefix, target.Os, target.Arch))
	if target.Os == consts.OsWindows {
		res += ".exe"
	}
	return res
}
