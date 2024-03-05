package docker

import "context"

func RemoveImage(ctx context.Context, image string) error {
	args := []string{"image", "rm", image}
	return Exec(ctx, args...)
}
