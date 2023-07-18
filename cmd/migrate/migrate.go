package migrate

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/sourcegraph/conc"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			keepRunning, err := cmd.Flags().GetBool("keep-running")
			if err != nil {
				return err
			}

			logger = logger.With("service", "migrate")

			if len(conf.Migration.Tasks) == 0 {
				return errors.New("no migration task was configured")
			}

			var undo bool
			var sources []migration.Source
			var migrators []migration.Migrator
			for _, t := range conf.Migration.Tasks {
				task := t
				source, err := Source(&task, logger)
				if err != nil {
					logger.Error(err)
					continue
				}

				if err := source.Connect(context.Background()); err != nil {
					logger.Errorf("task.connect: %v", err)
					continue
				}

				sources = append(sources, source)

				migrator, err := source.Migrator(t.Source)
				if err != nil {
					logger.Errorf("migrator.init: %v", err)
					continue
				}

				migrators = append(migrators, migrator)

				err = migrator.Up()
				if err != nil {
					logger.Errorf("migrator.up: %v", err)
					undo = true
					break
				}

				logger.Info("upped")
			}

			if undo && len(migrators) > 0 {
				var wg conc.WaitGroup
				for _, m := range migrators {
					migrator := m
					wg.Go(func() {
						if err := migrator.Down(); err != nil {
							logger.Errorf("migrator.dow: %v", err)
							return
						}

						logger.Info("downed")
					})
				}
				wg.Wait()
			}

			if len(sources) > 0 {
				var wg conc.WaitGroup
				for _, s := range sources {
					source := s
					wg.Go(func() {
						if err := source.Disconnect(context.Background()); err != nil {
							logger.Errorf("task.disconnect: %v", err)
						}
					})
				}
				wg.Wait()
			}

			if !keepRunning {
				return nil
			}

			if err := utils.Liveness("kanthor.migrate"); err != nil {
				return err
			}
			if err := utils.Readiness("kanthor.migrate"); err != nil {
				return err
			}

			ctx, _ := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			<-ctx.Done()
			return nil
		},
	}

	command.Flags().BoolP("keep-running", "", false, "--keep-running: force migration to run once finished to prevent it from keep restarting. It's useful when you deploy it on UAT/PROD.")

	return command
}