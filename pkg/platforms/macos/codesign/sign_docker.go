package codesign

import (
	"fmt"

	"github.com/cardinalby/xgo-pack/pkg/docker"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	fsutil "github.com/cardinalby/xgo-pack/pkg/util/fs"
)

const rcodesignImage = "ghcr.io/cardinalby/apple-codesign-docker:latest"
const containerWorkingDir = "/xgopack_root"

// SignArtifact signs the artifact (bin or app bundle) at the given path relative to the root dir.
// To customize the signing process, create rcodesign.toml in the root dir.
func SignArtifact(ctx buildctx.Context, artifactPath string) error {
	ctx.Logger.Printf("Signing '%s'...\n", artifactPath)
	artifactRelPath, err := fsutil.GetRelPathInsideDir(artifactPath, ctx.Cfg.Root)
	if err != nil {
		return err
	}

	volumes := map[string]string{
		ctx.Cfg.Root: containerWorkingDir,
	}

	stdout, err := docker.RunOnceRes(ctx, rcodesignImage, docker.RunOptions{
		Volumes:    volumes,
		WorkingDir: containerWorkingDir,
		Args:       []string{"sign", artifactRelPath},
	})
	if err != nil {
		return fmt.Errorf("error signing '%s': %w", artifactPath, err)
	}
	ctx.Logger.Println(stdout)
	return nil
}
