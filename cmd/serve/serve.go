package serve

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/scrapnode/kanthor/cmd/utils"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/internal/debugging"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
	"github.com/scrapnode/kanthor/services"
	"github.com/spf13/cobra"
)

var example = `
	# serve one service
	kanthor serve sdk
	# serve multiple services
	kanthor serve sdk portal schedule dispatcher storage
`

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:       "serve",
		Example:   example,
		ValidArgs: append(services.SERVICES, services.ALL),
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			if slices.Contains(args, services.ALL) || len(args) > 1 {
				return multiple(provider, args)
			}

			return single(provider, args[0])
		},
	}

	return command
}

func single(provider configuration.Provider, name string) error {
	logger, err := logging.New(provider)
	if err != nil {
		return err
	}

	service, err := Service(provider, name)
	if err != nil {
		return err
	}
	debug := debugging.NewServer()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = service.Start(ctx); err != nil {
		return err
	}
	if err = debug.Start(ctx); err != nil {
		return err
	}

	defer func() {
		if err := utils.Stop(service, debug); err != nil {
			logger.Error(err)
		}
	}()

	go func() {
		if err = service.Run(ctx); err != nil {
			logger.Error(err)
		}
		logger.Debug("exit running process")
	}()
	go func() {
		if err = debug.Run(ctx); err != nil {
			logger.Error(err)
		}
	}()

	// listen for the interrupt signal.
	<-ctx.Done()
	logger.Warnw("SYSTEM.SIGNAL.INTERRUPT", "error", ctx.Err(), "signal", fmt.Sprintf("%v", ctx))
	return nil
}

func multiple(provider configuration.Provider, names []string) error {
	logger, err := logging.New(provider)
	if err != nil {
		return err
	}

	instances, err := Services(provider, names)
	if err != nil {
		return err
	}

	defer func() {
		var items []utils.Stoppable
		for _, instance := range instances {
			items = append(items, instance.(utils.Stoppable))
		}
		if err := utils.Stop(items...); err != nil {
			logger.Error(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, instance := range instances {
		if err = instance.Start(ctx); err != nil {
			return err
		}
		go func(service patterns.Runnable) {
			if err = service.Run(ctx); err != nil {
				logger.Error(err)
			}
			logger.Debug("exit running process")
		}(instance)
	}

	// listen for the interrupt signal.
	<-ctx.Done()
	logger.Warnw("SYSTEM.SIGNAL.INTERRUPT", "error", ctx.Err())
	return nil
}
