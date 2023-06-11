package cmd

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/migration/ioc"
	"github.com/scrapnode/kanthor/migration/operators"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewUp() *cobra.Command {
	command := &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider, err := config.New()
			if err != nil {
				return err
			}

			logger, err := ioc.InitializeLogger(provider)
			if err != nil {
				return err
			}

			migrator, err := ioc.InitializeMigrator(provider)
			if err != nil {
				return err
			}

			keepRunning, err := cmd.Flags().GetBool("keep-running")
			if err != nil {
				return err
			}

			ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			if err := migrator.Connect(ctx); err != nil {
				return err
			}

			go func() {
				err := migrator.Up()

				if err != nil && !errors.Is(err, operators.ErrNoChange) {
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

			ctx, cancel = context.WithTimeout(cmd.Context(), 11*time.Second)
			go func() {
				if err := migrator.Disconnect(ctx); err != nil {
					logger.Error(err)
				}
				cancel()
			}()
			<-ctx.Done()
			return nil
		},
	}

	command.Flags().BoolP("keep-running", "", false, "--keep-running: keep migrator running after finished. It's useful when you deployed it on UAT/PROD")

	return command
}
