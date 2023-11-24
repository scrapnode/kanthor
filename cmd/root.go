package cmd

import (
	"github.com/scrapnode/kanthor/cmd/administer"
	"github.com/scrapnode/kanthor/cmd/check"
	"github.com/scrapnode/kanthor/cmd/config"
	"github.com/scrapnode/kanthor/cmd/migrate"
	"github.com/scrapnode/kanthor/cmd/serve"
	"github.com/scrapnode/kanthor/cmd/setup"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{}

	command.AddCommand(check.New())
	command.AddCommand(config.New(provider))
	command.AddCommand(migrate.New(provider))
	command.AddCommand(setup.New(provider))
	command.AddCommand(serve.New(provider))
	command.AddCommand(administer.New(provider))

	command.PersistentFlags().BoolP("verbose", "", false, "--verbose | show more information")
	return command
}
