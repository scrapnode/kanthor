package cmd

import (
	dataplane "github.com/scrapnode/kanthor/dataplane/cmd"
	"github.com/spf13/cobra"
)

func NewServe() *cobra.Command {
	command := &cobra.Command{
		Use: "serve",
	}
	command.AddCommand(dataplane.New())
	return command
}
