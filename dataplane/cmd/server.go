package cmd

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/dataplane/ioc"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewServer() *cobra.Command {
	command := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider, err := config.New()
			if err != nil {
				return err
			}

			logger, err := ioc.InitializeLogger(provider)
			if err != nil {
				return err
			}

			server, err := ioc.InitializeServer(provider)
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			if err := server.Start(ctx); err != nil {
				return err
			}

			go func() {
				if err := server.Run(ctx); err != nil {
					logger.Error(errors.Unwrap(err))
					cancel()
					return
				}
			}()

			// Listen for the interrupt signal.
			<-ctx.Done()
			// make sure once we stop process, we cancel all the execution
			cancel()

			ctx, cancel = context.WithTimeout(cmd.Context(), 11*time.Second)
			go func() {
				if err := server.Stop(ctx); err != nil {
					logger.Error(errors.Unwrap(err))
				}
				cancel()
			}()
			<-ctx.Done()
			return nil
		},
	}
	return command
}
