package docker

import (
	"context"
	"fmt"
	"os/user"
)

type RunOptions struct {
	Volumes map[string]string
	Env     map[string]string
	Args    []string
	// If not empty, sets uid/gid of the user inside the container. Otherwise, the user is the current user.
	User string
}

func RunOnce(ctx context.Context, image string, options RunOptions) error {
	args, err := prepareRunArgs(image, options)
	if err != nil {
		return err
	}
	return Exec(ctx, args...)
}

func RunOnceRes(ctx context.Context, image string, options RunOptions) (stdout string, err error) {
	args, err := prepareRunArgs(image, options)
	if err != nil {
		return "", err
	}
	return ExecRes(ctx, args...)
}

func prepareRunArgs(image string, options RunOptions) ([]string, error) {
	if options.User == "" {
		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("error getting current user: %w", err)
		}
		options.User = fmt.Sprintf("%s:%s", usr.Uid, usr.Gid)
	}
	args := append([]string{"run", "--rm"})
	for key, value := range options.Volumes {
		args = append(args, "-v", fmt.Sprintf(`%s:%s`, key, value))
	}
	for key, value := range options.Env {
		args = append(args, "-e", fmt.Sprintf(`%s=%s`, key, value))
	}
	args = append(args, "--user", options.User)
	args = append(args, image)
	args = append(args, options.Args...)

	return args, nil
}
