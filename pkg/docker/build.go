package docker

import (
	"context"
	"strings"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
)

type Image struct {
	isDisposed bool
	Sha        string
}

func (i *Image) Dispose() error {
	if !i.isDisposed {
		return RemoveImage(context.Background(), i.Sha)
	}
	return nil
}

func BuildDir(ctx buildctx.Context, dirPath string) (Image, error) {
	args := []string{"build", "-q", dirPath}
	stdout, err := ExecRes(ctx, args...)
	if err != nil {
		return Image{}, err
	}
	return Image{Sha: strings.TrimSpace(stdout)}, nil
}
