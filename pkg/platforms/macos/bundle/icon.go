package bundle

import (
	"fmt"
	"image"
	"os"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
)

func getPngImage(ctx buildctx.Context, pngImagePath string) (image.Image, error) {
	pngFile, err := os.Open(pngImagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if err := pngFile.Close(); err != nil {
			ctx.Logger.Printf("error closing '%s' file: %v", pngImagePath, err)
		}
	}()
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	srcImg, _, err := image.Decode(pngFile)
	if err != nil {
		return nil, fmt.Errorf("error decoding '%s' image: %w", pngImagePath, err)
	}
	return srcImg, nil
}
