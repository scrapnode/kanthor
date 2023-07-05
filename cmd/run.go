package cmd

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/ioc"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewRun(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:       "run",
		ValidArgs: []string{services.MIGRATION},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			if err := conf.Validate(serviceName); err != nil {
				return err
			}

			keepRunning, err := cmd.Flags().GetBool("keep-running")
			if err != nil {
				return err
			}

			runner, err := useRunner(serviceName, conf, logger)
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			if err := runner.Start(ctx); err != nil {
				return err
			}

			go func() {
				if err := runner.Run(ctx); err != nil {
					logger.Error(err)
					cancel()
					return
				}

				logger.Info("completed")

				if !keepRunning {
					cancel()
				}
			}()

			// Listen for the interrupt signal.
			<-ctx.Done()
			// make sure once we stop process, we cancel all the execution
			cancel()

			// wait a little to stop our service
			ctx, cancel = context.WithTimeout(cmd.Context(), 11*time.Second)
			go func() {
				if err := runner.Stop(ctx); err != nil {
					logger.Error(err)
				}
				cancel()
			}()
			<-ctx.Done()

			return nil
		},
	}

	command.Flags().BoolP("keep-running", "", false, "--keep-running: force migration running after finished to prevent it keep restarting. It's useful when you deployed it on UAT/PROD")

	return command
}

func useRunner(name string, conf *config.Config, logger logging.Logger) (services.Service, error) {
	if name == services.MIGRATION {
		return ioc.InitializeMigration(conf, logger)
	}

	return nil, fmt.Errorf("run.runner: unknow runner [%s]", name)
}
