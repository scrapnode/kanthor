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

	command.PersistentFlags().Int64P("timeout", "", 60000, "--timeout=600000 | default timeout in milliseconds")
	return command
}
