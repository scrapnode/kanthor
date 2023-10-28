package migrate

import (
	"errors"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/datastore/migrator"
	"github.com/scrapnode/kanthor/logging"
	"github.com/spf13/cobra"
)

func NewDatastore(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:       "datastore",
		ValidArgs: []string{"up", "down"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := logging.New(provider)
			if err != nil {
				return err
			}
			migrate, err := migrator.New(provider)
			if err != nil {
				return err
			}

			if args[0] == "up" {
				if err := migrate.Steps(1); err != nil {
					return err
				}
				version, dirty := migrate.Version()
				logger.Infow("up to next version", "version", version, "dirty", dirty)
				return nil
			}
			if args[0] == "down" {
				if err := migrate.Steps(-1); err != nil {
					return err
				}
				version, dirty := migrate.Version()
				logger.Infow("down to previous version", "version", version, "dirty", dirty)
				return nil
			}

			return errors.New("unknow migration target")
		},
	}

	return command
}
