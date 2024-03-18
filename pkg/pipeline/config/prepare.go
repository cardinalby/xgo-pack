package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cardinalby/go-dto-merge"
	"github.com/cardinalby/xgo-pack/pkg/go_src"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/presets"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/util"
	fsutil "github.com/cardinalby/xgo-pack/pkg/util/fs"
	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
)

func Prepare(c cfgtypes.Config, presetsResolver presets.Resolver) (cfg cfgtypes.Config, err error) {
	cfg, err = presets.ApplyConfigPreset(c, presetsResolver)
	if err != nil {
		return cfg, err
	}

	if err = FillDefaults(&cfg); err != nil {
		return cfg, fmt.Errorf("error filling defaults: %w", err)
	}

	if err = applyCommons(&cfg); err != nil {
		return cfg, fmt.Errorf("error merging common config for %w", err)
	}

	return cfg, nil
}

func FillDefaults(c *cfgtypes.Config) (err error) {
	if c.Schema == "" {
		c.Schema = cfgtypes.DefaultSchema
	}

	if c.Root == "" {
		c.Root, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("error obtaining working dir: %w", err)
		}
	}

	if c.DistDir == "" {
		c.DistDir, err = fsutil.GetNonExisting(c.Root, "dist")
		if err != nil {
			return fmt.Errorf("error obtaining default dist dir: %w", err)
		}
	}

	if c.TmpDir == "" {
		c.TmpDir, err = fsutil.GetNonExisting(c.Root, filepath.Join(c.DistDir, "tmp"))
		if err != nil {
			return fmt.Errorf("error obtaining default tmp dir: %w", err)
		}
	}

	if c.Src.MainPkg == "" {
		c.Src.MainPkg, err = go_src.FindMainPkgDir(c.Root)
		if err != nil {
			return fmt.Errorf("error looking for main package dir: %w", err)
		}
	}

	if c.Src.Icon == "" {
		c.Src.Icon = cfgtypes.DefaultBuiltInIcon
	}

	if c.Targets.Common.Identifier == "" {
		goModFilePath, err := go_src.FindGoModFile(c.Root)
		if err != nil {
			return fmt.Errorf("error looking for go.mod file: %w", err)
		}
		moduleName, err := go_src.GetModuleName(goModFilePath)
		if err != nil {
			return fmt.Errorf("error getting module name from '%s': %w", goModFilePath, err)
		}
		c.Targets.Common.Identifier = util.ModulePathToIdentifier(moduleName)
	}

	if c.Targets.Common.Copyright == "" {
		c.Targets.Common.Copyright = fmt.Sprintf(
			"© %d, %s",
			time.Now().Year(),
			util.IdentifierWithoutLastPart(c.Targets.Common.Identifier),
		)
	}

	if c.Targets.Common.ProductName == "" {
		c.Targets.Common.ProductName = util.IdentifierLastPart(c.Targets.Common.Identifier)
	}

	if c.Targets.Common.BinName == "" {
		mainPkgAbsPath, err := filepath.Abs(filepath.Join(c.Root, c.Src.MainPkg))
		if err != nil {
			return fmt.Errorf("error getting absolute path for main package '%s': %w", c.Src.MainPkg, err)
		}
		c.Targets.Common.BinName = filepath.Base(mainPkgAbsPath)
	}

	if c.Targets.Common.HighDpi == nil {
		c.Targets.Common.HighDpi = typeutil.Ptr(false)
	}

	if c.Targets.Common.Version == "" {
		c.Targets.Common.Version = "1.0.0"
	}

	for osName, osTargets := range c.Targets.GetOsArches() {
		for arch, archCfg := range osTargets.GetArches() {
			if archCfg.GetOutDir() == "" {
				archCfg.SetOutDir(fmt.Sprintf("%s_%s", osName, arch))
			}
		}
	}

	if c.Targets.Macos.Common.Codesign.Sign == nil {
		c.Targets.Macos.Common.Codesign.Sign = typeutil.Ptr(true)
	}

	if c.Targets.Macos.Common.Bundle.HideInDock == nil {
		c.Targets.Macos.Common.Bundle.HideInDock = typeutil.Ptr(false)
	}

	if c.Targets.Macos.Common.Bundle.BundleName == "" {
		c.Targets.Macos.Common.Bundle.BundleName = c.Targets.Common.ProductName + ".app"
	}

	if c.Targets.Macos.Common.Dmg.DmgName == "" {
		c.Targets.Macos.Common.Dmg.DmgName = c.Targets.Common.ProductName + ".dmg"
	}

	if c.Targets.Linux.Common.Deb.DebName == "" {
		c.Targets.Linux.Common.Deb.DebName = c.Targets.Common.ProductName + ".deb"
	}

	if c.Targets.Linux.Common.Deb.Name == "" {
		c.Targets.Linux.Common.Deb.Name = c.Targets.Common.ProductName
	}

	if c.Targets.Linux.Common.Deb.Section == "" {
		c.Targets.Linux.Common.Deb.Section = "default"
	}

	if c.Targets.Linux.Common.Deb.Maintainer == "" {
		c.Targets.Linux.Common.Deb.Maintainer = util.IdentifierWithoutLastPart(c.Targets.Common.Identifier)
	}

	if c.Targets.Linux.Common.Deb.DstBinPath == "" {
		c.Targets.Linux.Common.Deb.DstBinPath = "/usr/bin/" + c.Targets.Common.BinName
	}

	if c.Targets.Linux.Common.Deb.DesktopEntry.AddDesktopEntry == nil {
		c.Targets.Linux.Common.Deb.DesktopEntry.AddDesktopEntry = typeutil.Ptr(true)
	}

	if c.Targets.Linux.Common.Deb.DesktopEntry.AddIcon == nil {
		c.Targets.Linux.Common.Deb.DesktopEntry.AddIcon = typeutil.Ptr(true)
	}

	if c.Targets.Linux.Common.Deb.DesktopEntry.DstIconPath == "" {
		c.Targets.Linux.Common.Deb.DesktopEntry.DstIconPath = "/usr/share/icons/" + c.Targets.Common.Identifier + ".png"
	}

	if c.Targets.Linux.Common.Deb.DesktopEntry.Name == "" {
		c.Targets.Linux.Common.Deb.DesktopEntry.Name = c.Targets.Common.ProductName
	}

	if c.Targets.Linux.Common.Deb.DesktopEntry.Type == "" {
		c.Targets.Linux.Common.Deb.DesktopEntry.Type = "Application"
	}

	if c.Targets.Linux.Common.Deb.DesktopEntry.Terminal == nil {
		c.Targets.Linux.Common.Deb.DesktopEntry.Terminal = typeutil.Ptr(true)
	}

	if c.Targets.Linux.Common.Deb.DesktopEntry.NoDisplay == nil {
		c.Targets.Linux.Common.Deb.DesktopEntry.NoDisplay = typeutil.Ptr(false)
	}

	return nil
}

func applyCommons(c *cfgtypes.Config) (err error) {
	for osName, osTargets := range c.Targets.GetOsArches() {
		osCommonCfg, err := dtomerge.Merge(
			c.Targets.Common, osTargets.GetCommonCfg(),
		)
		if err != nil {
			return fmt.Errorf("%s: %w", osName, err)
		}
		osTargets.SetCommonCfg(osCommonCfg)
	}
	return nil
}
