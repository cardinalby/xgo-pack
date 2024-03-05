package docker

import (
	"context"
	"fmt"
)

type RunOptions struct {
	Volumes map[string]string
	Env     map[string]string
	Args    []string
}

func RunOnce(ctx context.Context, image string, options RunOptions) error {
	args := prepareRunArgs(image, options)
	return Exec(ctx, args...)
}

func RunOnceRes(ctx context.Context, image string, options RunOptions) (stdout string, err error) {
	args := prepareRunArgs(image, options)
	return ExecRes(ctx, args...)
}

func prepareRunArgs(image string, options RunOptions) []string {
	args := append([]string{"run", "--rm"})
	for key, value := range options.Volumes {
		args = append(args, "-v", fmt.Sprintf(`%s:%s`, key, value))
	}
	for key, value := range options.Env {
		args = append(args, "-e", fmt.Sprintf(`%s=%s`, key, value))
	}
	args = append(args, image)
	args = append(args, options.Args...)
	return args
}
