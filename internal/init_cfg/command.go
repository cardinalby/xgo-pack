package init_cfg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	typeutil "github.com/cardinalby/xgo-pack/pkg/util/type"
	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/cobra"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	fsutil "github.com/cardinalby/xgo-pack/pkg/util/fs"
)

func GetCommand(defaultCfgRelPath string) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialise config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialiseConfig(defaultCfgRelPath)
		},
	}
}

func initialiseConfig(defaultCfgRelPath string) error {
	// check if exists
	_, err := os.Stat(defaultCfgRelPath)
	if err == nil {
		if !askForConfirmation(fmt.Sprintf("Config file '%s' already exists. Overwrite?", defaultCfgRelPath)) {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking if config file '%s' exists: %w", defaultCfgRelPath, err)
	}

	var cfg cfgtypes.Config
	if err := config.FillDefaults(&cfg); err != nil {
		return err
	}
	cfg, commentMap := modifyInitCfg(cfg)
	data, err := yaml.MarshalWithOptions(cfg, yaml.WithComment(commentMap))
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := fsutil.WriteFile(defaultCfgRelPath, data); err != nil {
		return fmt.Errorf("error writing config file '%s': %w", defaultCfgRelPath, err)
	}

	fmt.Printf("Default config was saved to %s", defaultCfgRelPath)

	return nil
}

func modifyInitCfg(cfg cfgtypes.Config) (cfgtypes.Config, yaml.CommentMap) {
	comments := make(yaml.CommentMap)

	cfg.Root = ""
	comments["$.src.icon"] = []*yaml.Comment{yaml.LineComment(" Set path to your icon")}

	cfg.Targets.Windows.Amd64.BuildBin = typeutil.Ptr(true)
	comments["$.targets.windows.amd64.build_bin"] = []*yaml.Comment{yaml.LineComment(" Build and keep binary")}

	for arch, archCfg := range cfg.Targets.Linux.GetLinuxArches() {
		yamlPath := fmt.Sprintf("$.targets.linux.%s", arch)
		archCfg.BuildBin = typeutil.Ptr(true)
		comments[yamlPath+".build_bin"] = []*yaml.Comment{yaml.LineComment(" Build and keep binary")}
		archCfg.BuildDeb = typeutil.Ptr(true)
		comments[yamlPath+".build_deb"] = []*yaml.Comment{yaml.LineComment(" Build and keep deb package")}
	}
	for arch, archCfg := range cfg.Targets.Macos.GetMacosArches() {
		yamlPath := fmt.Sprintf("$.targets.macos.%s", arch)
		archCfg.BuildBin = typeutil.Ptr(true)
		comments[yamlPath+".build_bin"] = []*yaml.Comment{yaml.LineComment(" Build and keep binary")}
		archCfg.BuildBundle = typeutil.Ptr(true)
		comments[yamlPath+".build_bundle"] = []*yaml.Comment{yaml.LineComment(" Build and keep app bundle")}
		archCfg.BuildDmg = typeutil.Ptr(true)
		comments[yamlPath+".build_dmg"] = []*yaml.Comment{yaml.LineComment(" Build and keep dmg image with bundle")}
	}
	return cfg, comments
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
