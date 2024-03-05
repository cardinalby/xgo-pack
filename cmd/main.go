package main

import (
	"os"

	"github.com/cardinalby/xgo-pack/internal/init_cfg"
	"github.com/cardinalby/xgo-pack/internal/run_pipeline"
	"github.com/cardinalby/xgo-pack/pkg/util/logging"
	"github.com/spf13/cobra"
)

const defaultCfgRelPath = "xgo-pack-config.yaml"

func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(
		run_pipeline.GetCommand(defaultCfgRelPath),
		init_cfg.GetCommand(defaultCfgRelPath),
	)

	errLogger := logging.NewLogger(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		errLogger.Println(err.Error())
		os.Exit(1)
	}
}
