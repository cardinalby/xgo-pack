package buildctx

import (
	"fmt"
	"strconv"
	"strings"
)

const noLayerIndex = -1

// ImgSourcePath points to a source image for conversion. Can be just a file path or can have optional
// format specifier and layer index:
// FORMAT:some/file/path.png[layerIndex]
// where layerIndex is a positive index of a layer in a multi-layer image
type ImgSourcePath string
type ImgSourcePathComponents struct {
	Format     string
	FilePath   string
	LayerIndex int
}

func (s ImgSourcePathComponents) ToSourcePath() ImgSourcePath {
	sb := strings.Builder{}
	if s.Format != "" {
		sb.WriteString(s.Format)
		sb.WriteString(":")
	}
	sb.WriteString(s.FilePath)
	if s.LayerIndex != noLayerIndex {
		sb.WriteString("[")
		sb.WriteString(strconv.Itoa(s.LayerIndex))
		sb.WriteString("]")
	}
	return ImgSourcePath(sb.String())
}

func (s ImgSourcePath) Components() (components ImgSourcePathComponents, err error) {
	parts := strings.SplitN(string(s), ":", 2)
	if len(parts) == 1 {
		components.Format = ""
		components.FilePath = parts[0]
	} else {
		components.Format = parts[0]
		components.FilePath = parts[1]
	}
	parts = strings.Split(components.FilePath, "[")
	if len(parts) == 1 {
		components.LayerIndex = noLayerIndex
		return
	}
	if len(parts) != 2 {
		err = fmt.Errorf("multiple '[' sign occurence")
		return
	}
	components.FilePath = parts[0]
	components.LayerIndex, err = strconv.Atoi(strings.TrimSuffix(parts[1], "]"))
	if err != nil || components.LayerIndex < 0 {
		err = fmt.Errorf("invalid layer index in source path: %s", s)
		return
	}
	return
}
