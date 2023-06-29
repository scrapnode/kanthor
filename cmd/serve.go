package cmd

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func NewServe(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "serve",
	}
	command.AddCommand(NewServeDataplane(conf, logger))
	command.AddCommand(NewServeScheduler(conf, logger))
	return command
}
