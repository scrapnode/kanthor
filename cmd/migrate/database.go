package migrate

import (
	"errors"

	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/spf13/cobra"
)

func NewDatabase(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:       "database",
		ValidArgs: []string{"up", "down"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			instance, err := database.New(&conf.Database, logger)
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

			m, err := instance.Migrator()
			if err != nil {
				return err
			}

			if args[0] == "up" {
				if err := m.Steps(1); err != nil {
					return err
				}
				version, dirty := m.Version()
				logger.Infow("up to next version", "version", version, "dirty", dirty)
				return nil
			}
			if args[0] == "down" {
				if err := m.Steps(-1); err != nil {
					return err
				}
				version, dirty := m.Version()
				logger.Infow("down to previous version", "version", version, "dirty", dirty)
				return nil
			}

			return errors.New("unknow migration target")
		},
	}

	return command
}
