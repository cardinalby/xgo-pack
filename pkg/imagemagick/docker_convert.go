package imagemagick

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cardinalby/xgo-pack/pkg/docker"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
)

// runConvert runs docker container with imagemagick convert command.
// srcPath - source image path
// dstPath - destination image path
// getArgsClb - callback to get convert command arguments
func runConvert(
	ctx context.Context,
	srcPath buildctx.ImgSourcePath,
	dstPath string,
	getArgsClb func(src, dst string) []string,
) error {
	srcComponents, err := srcPath.Components()
	if err != nil {
		return fmt.Errorf("invalid source path '%s': %w", srcPath, err)
	}

	srcAbsPath, err := filepath.Abs(srcComponents.FilePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path for '%s': %w", srcComponents.FilePath, err)
	}

	dstAbsPath, err := filepath.Abs(dstPath)
	if err != nil {
		return fmt.Errorf("error getting absolute path for '%s': %w", dstPath, err)
	}

	volumes := map[string]string{
		srcAbsPath:               "/src_f",
		filepath.Dir(dstAbsPath): "/dst",
	}

	srcComponents.FilePath = "/src_f"
	args := getArgsClb(
		string(srcComponents.ToSourcePath()),
		"/dst/"+filepath.Base(dstPath),
	)

	return docker.RunOnce(ctx, dockerImage, docker.RunOptions{
		Volumes: volumes,
		Args:    args,
	})
}
