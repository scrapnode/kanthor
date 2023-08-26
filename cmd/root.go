package cmd

import (
	"github.com/scrapnode/kanthor/cmd/check"
	"github.com/scrapnode/kanthor/cmd/do"
	"github.com/scrapnode/kanthor/cmd/migrate"
	"github.com/scrapnode/kanthor/cmd/serve"
	"github.com/scrapnode/kanthor/cmd/show"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider, conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{}

	command.AddCommand(NewVersion(provider, conf))
	command.AddCommand(do.New(conf, logger))
	command.AddCommand(show.New(provider, conf))
	command.AddCommand(migrate.New(conf, logger))
	command.AddCommand(serve.New(conf, logger))
	command.AddCommand(check.New())

	command.PersistentFlags().BoolP("verbose", "", false, "--verbose: show more information")
	return command
}
