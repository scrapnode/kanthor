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
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:       "serve",
		ValidArgs: append(services.SERVICES, services.ALL),
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if name != services.ALL {
				return single(provider, name)
			}

			return multiple(provider)
		},
	}

	return command
}

func single(provider configuration.Provider, name string) (err error) {
	logger, err := logging.New(provider)
	if err != nil {
		return
	}

	service, err := Service(provider, name)
	if err != nil {
		return
	}
	debug := debugging.NewServer()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err = service.Start(ctx); err != nil {
		return
	}
	if err = debug.Start(ctx); err != nil {
		return
	}

	go func() {
		if err = service.Run(ctx); err != nil {
			logger.Error(err)
		}
	}()
	go func() {
		if err = debug.Run(ctx); err != nil {
			logger.Error(err)
		}
	}()

	defer func() {
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
		case err = <-errc:
			return
		case <-ctx.Done():
			err = ctx.Err()
		}
	}()

	// listen for the interrupt signal.
	<-ctx.Done()
	logger.Warnw("SYSTEM.SIGNAL.INTERRUPT", "error", ctx.Err())
	return
}

func multiple(provider configuration.Provider) (err error) {
	logger, err := logging.New(provider)
	if err != nil {
		return
	}

	instances, err := Services(provider)
	if err != nil {
		return
	}

	defer func() {
		// wait a little to stop our service
		errc := make(chan error)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		go func() {
			var returning error
			for _, instance := range instances {
				if err := instance.Stop(ctx); err != nil {
					returning = errors.Join(returning, err)
				}
			}
			errc <- returning
		}()

		select {
		case err = <-errc:
			return
		case <-ctx.Done():
			err = ctx.Err()
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	for _, instance := range instances {
		if err = instance.Start(ctx); err != nil {
			return
		}
		go func(service patterns.Runnable) {
			if err = service.Run(ctx); err != nil {
				logger.Error(err)
			}
		}(instance)
	}

	// listen for the interrupt signal.
	<-ctx.Done()
	logger.Warnw("SYSTEM.SIGNAL.INTERRUPT", "error", ctx.Err())
	return
}
