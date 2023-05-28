package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "dataplane",
	}
	command.AddCommand(NewServer())
	return command
}
