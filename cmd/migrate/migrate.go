package migrate

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/sourcegraph/conc"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			keepRunning, err := cmd.Flags().GetBool("keep-running")
			if err != nil {
				return err
			}

			if len(conf.Migration.Tasks) == 0 {
				return errors.New("no migration task was configured")
			}

			var undo bool
			sources := []migration.Source{}
			migrators := []migration.Migrator{}
			for _, t := range conf.Migration.Tasks {
				mlogger := logger.With("service", "migrate", "task", t.Name)

				task := t
				source, err := Source(&task, mlogger)
				if err != nil {
					mlogger.Error(err)
					continue
				}

				if err := source.Connect(context.Background()); err != nil {
					mlogger.Errorf("task.connect: %v", err)
					continue
				}

				sources = append(sources, source)

				migrator, err := source.Migrator(t.Source)
				if err != nil {
					mlogger.Errorf("migrator.init: %v", err)
					continue
				}

				migrators = append(migrators, migrator)

				err = migrator.Up()
				if err != nil {
					mlogger.Errorf("migrator.up: %v", err)
					undo = true
					break
				}

				mlogger.Info("upped")
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

			ctx, _ := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			<-ctx.Done()
			return nil
		},
	}

	command.Flags().BoolP("keep-running", "", false, "--keep-running: force migration to run once finished to prevent it from keep restarting. It's useful when you deploy it on UAT/PROD.")

	return command
}
