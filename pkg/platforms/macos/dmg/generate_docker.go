package dmg

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/docker"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
)

// image that doesn't support creating "Applications" symlink
const createDmgDockerImage = "ghcr.io/cardinalby/create-dmg:latest"

func Generate(
	ctx buildctx.Context,
	srcBundlePath string,
	outDmgPath string,
	appName string,
	addApplicationsSymlink bool,
) error {
	internalSrcPath := "/src"
	internalSrcBundlePath := path.Join(internalSrcPath, appName+".app")
	internalDstDir := "/out"
	var env map[string]string
	if addApplicationsSymlink {
		env = map[string]string{
			"APPLICATIONS_SYMLINK": "1",
		}
	}

	srcBundlePathAbs, err := filepath.Abs(srcBundlePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path for src bundle '%s': %w", srcBundlePath, err)
	}

	outDmgPathDir := filepath.Dir(outDmgPath)
	outDmgDirPathAbs, err := filepath.Abs(outDmgPathDir)
	if err != nil {
		return fmt.Errorf("error getting absolute path for out dmg dir '%s': %w", outDmgPathDir, err)
	}

	stdout, err := docker.RunOnceRes(ctx, createDmgDockerImage, docker.RunOptions{
		Volumes: map[string]string{
			srcBundlePathAbs: internalSrcBundlePath,
			outDmgDirPathAbs: internalDstDir,
		},
		Env:  env,
		Args: []string{appName, internalSrcPath, path.Join(internalDstDir, filepath.Base(outDmgPath))},
	})
	if err != nil {
		return fmt.Errorf("error running create-dmg docker container: %w. %s", err, stdout)
	}
	ctx.Logger.Printf(stdout)
	return nil
}
