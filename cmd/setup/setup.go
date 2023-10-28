package setup

import (
	"github.com/scrapnode/kanthor/cmd/setup/account"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use: "setup",
	}

	command.AddCommand(account.New(provider))

	return command
}
