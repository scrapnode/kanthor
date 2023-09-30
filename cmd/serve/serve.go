package serve

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scrapnode/kanthor/cmd/show"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "serve",
		ValidArgs: []string{
			services.PORTAL_API,
			services.SDK_API,
			services.SCHEDULER,
			services.DISPATCHER,
			services.STORAGE,
		},
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			verbose, _ := cmd.Flags().GetBool("verbose")

			if err := conf.Validate(serviceName); err != nil {
				_ = show.Config(conf, []configuration.Source{}, verbose, false)
				return err
			}

			service, err := Service(serviceName, conf, logger)
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			if err := service.Start(ctx); err != nil {
				return err
			}

			go func() {
				if err := service.Run(ctx); err != nil {
					logger.Error(err)
					stop()
					return
				}

				stop()
			}()

			// listen for the interrupt signal.
			<-ctx.Done()

			// wait a little to stop our service
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
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
