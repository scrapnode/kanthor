package create

import "github.com/spf13/cobra"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "create",
	}

	command.AddCommand(Message())

	return command
}
