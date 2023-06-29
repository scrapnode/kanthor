package cmd

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func NewRun(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "run",
	}
	command.AddCommand(NewRunMigration(conf, logger))
	return command
}
