package cmd

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
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
		ValidArgs: []string{services.DATAPLANE, services.SCHEDULER, services.DISPATCHER},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			verbose, _ := cmd.Flags().GetBool("verbose")

			if err := conf.Validate(serviceName); err != nil {
				if verbose {
					// if we got any error, should show the current configuration to easier debugging
					_ = showConfig(conf, []configuration.Source{}, false, false)
				}
				return err
			}

			service, err := useService(serviceName, conf, logger)
			if err != nil {
				return err
			}

			exporter, err := useMetricExporter(serviceName, conf, logger)
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			if err := service.Start(ctx); err != nil {
				return err
			}

			if err := exporter.Start(ctx); err != nil {
				return err
			}

			go func() {
				if err := service.Run(ctx); err != nil {
					logger.Error(err)
					cancel()
					return
				}
			}()

			go func() {
				if err := exporter.Run(ctx); err != nil {
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

				if err := exporter.Stop(ctx); err != nil {
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
	if name == services.DATAPLANE {
		return ioc.InitializeDataplane(conf, logger)
	}
	if name == services.SCHEDULER {
		return ioc.InitializeScheduler(conf, logger)
	}
	if name == services.DISPATCHER {
		return ioc.InitializeDispatcher(conf, logger)
	}

	return nil, fmt.Errorf("serve.service: unknow service [%s]", name)
}

func useMetricExporter(name string, conf *config.Config, logger logging.Logger) (patterns.Runnable, error) {
	if name == services.DATAPLANE {
		return metric.NewExporter(&conf.Dataplane.Metrics, logger), nil
	}
	if name == services.SCHEDULER {
		return metric.NewExporter(&conf.Scheduler.Metrics, logger), nil
	}
	if name == services.DISPATCHER {
		return metric.NewExporter(&conf.Dispatcher.Metrics, logger), nil
	}

	return nil, fmt.Errorf("serve.metric.exporter: unknow service [%s]", name)
}
