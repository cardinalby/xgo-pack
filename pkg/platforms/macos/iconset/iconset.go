package iconset

import (
	"bytes"
	"fmt"
	"image"

	"github.com/jackmordaunt/icns"
)

func Generate(srcIcon image.Image) ([]byte, error) {
	var dstBuff bytes.Buffer
	if err := icns.Encode(&dstBuff, srcIcon); err != nil {
		return nil, fmt.Errorf("error encoding icns: %w", err)
	}
	return dstBuff.Bytes(), nil
}
