package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "migration",
	}

	command.AddCommand(NewUp())
	return command
}
