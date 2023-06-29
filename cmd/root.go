package cmd

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider, conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{}

	command.AddCommand(NewShow(provider, conf))
	command.AddCommand(NewServe(conf, logger))
	command.AddCommand(NewRun(conf, logger))
	return command
}
