package cmd

import (
	"github.com/scrapnode/kanthor/cmd/check"
	"github.com/scrapnode/kanthor/cmd/client"
	"github.com/scrapnode/kanthor/cmd/config"
	"github.com/scrapnode/kanthor/cmd/migrate"
	"github.com/scrapnode/kanthor/cmd/serve"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{}

	command.AddCommand(check.New())
	command.AddCommand(config.New(provider))
	command.AddCommand(migrate.New(provider))
	command.AddCommand(serve.New(provider))
	command.AddCommand(client.New())

	command.PersistentFlags().BoolP("verbose", "", false, "--verbose | show more information")
	return command
}
