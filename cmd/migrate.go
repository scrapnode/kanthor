package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/sourcegraph/conc"
	"github.com/spf13/cobra"
)

func NewMigrate(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(conf.Migration.Tasks) == 0 {
				return errors.New("no migration task was configured")
			}

			var undo bool
			var tasks []patterns.Migrator
			var migrators []patterns.Migrate
			for _, t := range conf.Migration.Tasks {
				task, err := useMigrationTask(&t, logger)
				if err != nil {
					logger.Error(err)
					continue
				}

				if err := task.Connect(context.Background()); err != nil {
					logger.Errorf("task.connect: %v", err)
					continue
				}

				tasks = append(tasks, task)

				migrator, err := task.Migrator(t.Source)
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

			if len(tasks) > 0 {
				var wg conc.WaitGroup
				for _, t := range tasks {
					task := t
					wg.Go(func() {
						if err := task.Disconnect(context.Background()); err != nil {
							logger.Errorf("task.disconnect: %v", err)
						}
					})
				}
				wg.Wait()
			}

			return nil
		},
	}

	command.Flags().BoolP("keep-running", "", false, "--keep-running: force migration running after finished to prevent it keep restarting. It's useful when you deployed it on UAT/PROD")

	return command
}

func useMigrationTask(task *config.MigrationTask, logger logging.Logger) (patterns.Migrator, error) {
	if task.Name == "database" {
		return database.New(&database.Config{Uri: task.Uri}, logger), nil
	}
	if task.Name == "datastore" {
		return datastore.New(&datastore.Config{Uri: task.Uri}, logger), nil
	}

	return nil, fmt.Errorf("migrate: unsupport task [%s]", task.Name)
}
