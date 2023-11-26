package administer

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use: "administer",
	}

	command.AddCommand(NewAttemptTrigger(provider))
	command.AddCommand(NewAttemptEndeavor(provider))

	return command
}
