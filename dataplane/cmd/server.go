package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

func NewServer() *cobra.Command {
	command := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no implementation")
		},
	}
	return command
}
