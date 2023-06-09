package cmd

import (
	migration "github.com/scrapnode/kanthor/migration/cmd"
	"github.com/spf13/cobra"
)

func NewRun() *cobra.Command {
	command := &cobra.Command{
		Use: "run",
	}
	command.AddCommand(migration.New())
	return command
}
