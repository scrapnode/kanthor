package migrate

import (
	"errors"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/migration"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:       "migrate",
		ValidArgs: []string{"database", "datastore", "up", "down"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(2), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			var m migration.Migrator
			if args[0] == "database" {
				instance, err := database.New(&conf.Database, logger, timer.New())
				if err != nil {
					return err
				}

				if err := instance.Connect(ctx); err != nil {
					return err
				}
				defer func() {
					if err := instance.Disconnect(ctx); err != nil {
						logger.Error(err)
					}
				}()

				m, err = instance.Migrator()
				if err != nil {
					return err
				}
			}

			if args[0] == "datastore" {
				instance, err := datastore.New(&conf.Datastore, logger, timer.New())
				if err != nil {
					return err
				}

				if err := instance.Connect(ctx); err != nil {
					return err
				}
				defer func() {
					if err := instance.Disconnect(ctx); err != nil {
						logger.Error(err)
					}
				}()

				m, err = instance.Migrator()
				if err != nil {
					return err
				}
			}

			if args[1] == "up" {
				if err := m.Steps(1); err != nil {
					return err
				}
				logger.Info("up to next version")
				return nil
			}
			if args[1] == "down" {
				if err := m.Steps(-1); err != nil {
					return err
				}
				logger.Info("down to previous version")
				return nil
			}

			return errors.New("unknow migration target")
		},
	}

	return command
}
