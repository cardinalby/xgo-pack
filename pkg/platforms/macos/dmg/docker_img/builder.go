package docker_img

import (
	"fmt"

	"github.com/cardinalby/xgo-pack/pkg/docker"
	"github.com/cardinalby/xgo-pack/pkg/docker/docker_resource"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	"github.com/cardinalby/xgo-pack/pkg/util/fs"
)

func BuildCreateDmgDockerImage(ctx buildctx.Context) (*docker_resource.TempImage, error) {
	dockerTempDir, err := ctx.NewTempDir("macos-dmg")
	if err != nil {
		return nil, fmt.Errorf("error creating temp dir for docker image: %w", err)
	}
	defer func() {
		if err := dockerTempDir.Dispose(); err != nil {
			ctx.Logger.Printf("Failed to remove temp dir '%s': %v", dockerTempDir.Path, err)
		}
	}()

	if err := fsutil.CopyFs(imageFS, dockerTempDir.GetPath()); err != nil {
		return nil, fmt.Errorf("error copying embed docker image fs to temp dir: %w", err)
	}

	dockerImage, err := docker.BuildDir(ctx, dockerTempDir.GetPath())
	if err != nil {
		return nil, fmt.Errorf("error building docker image: %w", err)
	}

	return docker_resource.NewTempImage(dockerImage.Sha), nil
}
