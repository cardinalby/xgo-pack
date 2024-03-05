package docker_img

import "embed"

//go:embed Dockerfile
//go:embed create-dmg-with-app_symlink.sh
var imageFS embed.FS
