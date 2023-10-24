package serve

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scrapnode/kanthor/cmd/show"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/debugging"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:       "serve",
		ValidArgs: config.SERVICES,
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
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
			debug := debugging.NewServer()
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
