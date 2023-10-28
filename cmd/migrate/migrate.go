package migrate

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "migrate data",
	}

	command.AddCommand(NewDatabase(provider))
	command.AddCommand(NewDatastore(provider))

	return command
}
