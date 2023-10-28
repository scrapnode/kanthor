package cmd

import (
	"github.com/scrapnode/kanthor/cmd/check"
	"github.com/scrapnode/kanthor/cmd/migrate"
	"github.com/scrapnode/kanthor/cmd/serve"
	"github.com/scrapnode/kanthor/cmd/setup"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{}

	command.AddCommand(NewVersion())
	command.AddCommand(check.New())
	command.AddCommand(migrate.New(provider))
	command.AddCommand(setup.New(conf, logger))
	command.AddCommand(serve.New(conf, logger))

	command.PersistentFlags().BoolP("verbose", "", false, "--verbose: show more information")
	return command
}
