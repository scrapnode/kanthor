package serve

import (
	"context"
	"github.com/scrapnode/kanthor/cmd/show"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:       "serve",
		ValidArgs: []string{services.PORTAL, services.DATAPLANE, services.SCHEDULER, services.DISPATCHER},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			verbose, _ := cmd.Flags().GetBool("verbose")

			if err := conf.Validate(serviceName); err != nil {
				if verbose {
					// if we got any error, should show the current configuration for easier debugging
					_ = show.Config(conf, []configuration.Source{}, false, false)
				}
				return err
			}

			service, err := Service(serviceName, conf, logger)
			if err != nil {
				return err
			}

			exporter, err := MetricExporter(serviceName, conf, logger)
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			if err := service.Start(ctx); err != nil {
				return err
			}
			if err := exporter.Start(ctx); err != nil {
				return err
			}

			go func() {
				if err := service.Run(ctx); err != nil {
					logger.Error(err)
					stop()
					return
				}
			}()

			go func() {
				if err := exporter.Run(ctx); err != nil {
					logger.Error(err)
					stop()
					return
				}
			}()

			// listen for the interrupt signal.
			<-ctx.Done()
			// restore default behavior on the interrupt signal and notify user of shutdown.
			stop()

			// wait a little to stop our service
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			go func() {
				if err := exporter.Stop(ctx); err != nil {
					logger.Error(err)
				}

				if err := service.Stop(ctx); err != nil {
					logger.Error(err)
				}
			}()
			<-ctx.Done()

			return nil
		},
	}

	return command
}
