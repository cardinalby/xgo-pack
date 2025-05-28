package docker_resource

import (
	"context"

	"github.com/cardinalby/xgo-pack/pkg/docker"
)

type TempImage struct {
	IsDisposed bool
	ImageID    string
}

func NewTempImage(imageID string) *TempImage {
	return &TempImage{
		ImageID: imageID,
	}
}

func (t *TempImage) Dispose() error {
	if t != nil && !t.IsDisposed {
		return docker.RemoveImage(context.Background(), t.ImageID)
	}
	return nil
}

func (t *TempImage) GetPath() string {
	return t.ImageID
}
