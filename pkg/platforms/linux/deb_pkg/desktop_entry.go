package deb_pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/buildctx"
	fsutil "github.com/cardinalby/xgo-pack/pkg/util/fs"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

type DesktopEntry struct {
	Name      string
	Type      string
	Exec      string
	Terminal  bool
	Icon      string
	NoDisplay bool
}

func (d *DesktopEntry) Marshall() []byte {
	sb := strings.Builder{}
	sb.WriteString("[Desktop Entry]\n")
	sb.WriteString("Name=" + d.Name + "\n")
	sb.WriteString("Type=" + d.Type + "\n")
	sb.WriteString("Exec=" + d.Exec + "\n")
	sb.WriteString("Terminal=" + strconv.FormatBool(d.Terminal) + "\n")
	sb.WriteString("Icon=" + d.Icon + "\n")
	sb.WriteString("NoDisplay=" + strconv.FormatBool(d.NoDisplay) + "\n")
	return []byte(sb.String())
}

func registerDesktopEntryBuilder(ctx buildctx.Context) {
	ctx.Artifacts.RegisterBuilder(buildctx.KindLinuxDesktopEntry, func(ctx buildctx.Context) (buildctx.Artifact, error) {
		file, err := ctx.NewTempFile("linux-desktop-entry.desktop")
		if err != nil {
			return nil, fmt.Errorf("error creating temp desktop entry file path: %w", err)
		}

		entry := DesktopEntry{
			Name:      ctx.Cfg.Targets.Linux.Common.Deb.DesktopEntry.Name,
			Type:      ctx.Cfg.Targets.Linux.Common.Deb.DesktopEntry.Type,
			Exec:      ctx.Cfg.Targets.Linux.Common.Deb.DstBinPath,
			Terminal:  typeutil.PtrValueOr(ctx.Cfg.Targets.Linux.Common.Deb.DesktopEntry.Terminal, true),
			Icon:      ctx.Cfg.Targets.Linux.Common.Deb.DesktopEntry.DstIconPath,
			NoDisplay: typeutil.PtrValueOrDefault(ctx.Cfg.Targets.Linux.Common.Deb.DesktopEntry.NoDisplay),
		}
		if err := fsutil.WriteFile(file.Path, entry.Marshall()); err != nil {
			return nil, fmt.Errorf("error writing '%s' file: %w", file.Path, err)
		}
		return file, nil
	})
}
