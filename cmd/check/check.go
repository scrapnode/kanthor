package check

import "github.com/spf13/cobra"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "check",
	}
	command.AddCommand(NewReadiness())
	command.AddCommand(NewLiveness())
	return command
}
