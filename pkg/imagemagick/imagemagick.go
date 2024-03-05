package imagemagick

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
)

const dockerImage = "dpokidov/imagemagick"
const icoFormat = "ICO"
const pngFormat = "PNG"

func ToPng(ctx context.Context, srcPath buildctx.ImgSourcePath, dstPath string) error {
	return runConvert(ctx, srcPath, dstPath, func(src, dst string) []string {
		return []string{
			src,
			fmt.Sprintf("%s:%s", pngFormat, dst),
		}
	})
}

func ToIco(ctx context.Context, srcPath buildctx.ImgSourcePath, dstPath string, icoSizes []int) error {
	var sizesStrArr []string
	for _, size := range icoSizes {
		sizesStrArr = append(sizesStrArr, strconv.Itoa(size))
	}
	sizesStr := strings.Join(sizesStrArr, ",")

	return runConvert(ctx, srcPath, dstPath, func(src, dst string) []string {
		return []string{
			src,
			"-define",
			"icon:auto-resize=" + sizesStr,
			fmt.Sprintf("%s:%s", icoFormat, dst),
		}
	})
}
