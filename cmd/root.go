package cmd

import (
	"github.com/scrapnode/kanthor/cmd/migrate"
	"github.com/scrapnode/kanthor/cmd/serve"
	"github.com/scrapnode/kanthor/cmd/setup"
	"github.com/scrapnode/kanthor/cmd/show"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider, conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{}

	command.AddCommand(NewVersion(provider, conf))
	command.AddCommand(show.New(provider, conf))
	command.AddCommand(setup.New(conf, logger))
	command.AddCommand(migrate.New(conf, logger))
	command.AddCommand(serve.New(conf, logger))

	command.PersistentFlags().BoolP("verbose", "", false, "--verbose: show more information")
	return command
}
