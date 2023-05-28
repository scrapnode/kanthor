package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	command := &cobra.Command{}
	command.AddCommand(NewServe())
	return command
}
