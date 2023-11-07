package serve

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/debugging"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:       "serve",
		ValidArgs: services.SERVICES,
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := logging.New(provider)
			if err != nil {
				return err
			}

			serviceName := args[0]
			service, err := Service(serviceName, provider)
			if err != nil {
				return err
			}
			debug := debugging.NewServer()

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			if err := service.Start(ctx); err != nil {
				return err
			}
			if err := debug.Start(ctx); err != nil {
				return err
			}

			go func() {
				if err := service.Run(ctx); err != nil {
					logger.Error(err)
				}
			}()
			go func() {
				if err := debug.Run(ctx); err != nil {
					logger.Error(err)
				}
			}()

			// listen for the interrupt signal.
			<-ctx.Done()

			logger.Warnw("got stop signal", "error", ctx.Err())

			// wait a little to stop our service
			errc := make(chan error)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			go func() {
				var returning error
				if err := service.Stop(ctx); err != nil {
					returning = errors.Join(returning, err)
				}
				if err := debug.Stop(ctx); err != nil {
					returning = errors.Join(returning, err)
				}

				errc <- returning
			}()

			select {
			case err := <-errc:
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}

	return command
}
