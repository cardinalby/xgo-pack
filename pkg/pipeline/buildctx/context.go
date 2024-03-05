package buildctx

import (
	"context"
	"time"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	"github.com/cardinalby/xgo-pack/pkg/util/fs/fs_resource"
	"github.com/cardinalby/xgo-pack/pkg/util/logging"
)

type Context struct {
	TempDir   string
	Logger    logging.Logger
	Ctx       context.Context
	Artifacts ArtifactsI
	Cfg       cfgtypes.Config
}

func NewContext(ctx context.Context, cfg cfgtypes.Config, artifacts ArtifactsI, logger logging.Logger) Context {
	return Context{
		TempDir:   cfg.GetTmpDirPath(),
		Logger:    logger,
		Ctx:       ctx,
		Cfg:       cfg,
		Artifacts: artifacts,
	}
}

func (c Context) WithContext(ctx context.Context) Context {
	c.Ctx = ctx
	return c
}

func (c Context) WriteTempFile(bytes []byte, fileName string) (*fs_resource.TempFile, error) {
	f, err := fs_resource.WriteTempFile(bytes, c.TempDir, fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c Context) NewTempFile(fileName string) (*fs_resource.TempFile, error) {
	return fs_resource.NewTempFile(c.TempDir, fileName)
}

func (c Context) NewTempDir(dirName string) (*fs_resource.TempDir, error) {
	return fs_resource.NewTempDir(c.TempDir, dirName)
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.Ctx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.Ctx.Done()
}

func (c Context) Err() error {
	return c.Ctx.Err()
}

func (c Context) Value(key any) any {
	return c.Ctx.Value(key)
}
