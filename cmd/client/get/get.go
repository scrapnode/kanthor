package get

import "github.com/spf13/cobra"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "get",
	}

	command.AddCommand(Account())

	return command
}
