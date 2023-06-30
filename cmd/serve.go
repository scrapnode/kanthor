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

func NewServe(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:       "serve",
		ValidArgs: []string{"dataplane", "scheduler", "dispatcher"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := useService(args[0], conf, logger)
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			if err := service.Start(ctx); err != nil {
				return err
			}

			go func() {
				if err := service.Run(ctx); err != nil {
					logger.Error(err)
					cancel()
					return
				}
			}()

			// Listen for the interrupt signal.
			<-ctx.Done()
			// make sure once we stop process, we cancel all the execution
			cancel()

			// wait a little to stop our service
			ctx, cancel = context.WithTimeout(cmd.Context(), 11*time.Second)
			go func() {
				if err := service.Stop(ctx); err != nil {
					logger.Error(err)
				}
				cancel()
			}()
			<-ctx.Done()

			return nil
		},
	}

	return command
}

func useService(name string, conf *config.Config, logger logging.Logger) (services.Service, error) {
	if name == "dataplane" {
		return ioc.InitializeDataplane(conf, logger)
	}
	if name == "scheduler" {
		return ioc.InitializeScheduler(conf, logger)
	}
	if name == "dispatcher" {
		return ioc.InitializeDispatcher(conf, logger)
	}

	return nil, fmt.Errorf("serve: unknow service [%s]", name)
}
