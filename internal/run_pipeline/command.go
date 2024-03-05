package run_pipeline

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/cardinalby/xgo-pack/pkg/pipeline"
	"github.com/cardinalby/xgo-pack/pkg/pipeline/config"
	"github.com/cardinalby/xgo-pack/pkg/util/logging"
	"github.com/spf13/cobra"
)

func GetCommand(defaultCfgRelPath string) *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Initialise config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			a := getArgs(args, defaultCfgRelPath)
			logger := logging.NewNopLogger()
			if !a.Quiet {
				logger = logging.NewLogger(os.Stdout)
			}

			cfg, err := config.ReadFile(filepath.Join(a.Root, a.Cfg), a.Root)
			if err != nil {
				return fmt.Errorf("error reading config file '%s': %w", a.Cfg, err)
			}
			return pipeline.Start(getAppContext(), logger, cfg)
		},
	}
}

func getAppContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	call := make(chan os.Signal, 1)
	signal.Notify(call, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-call
		cancel()
	}()
	return ctx
}
