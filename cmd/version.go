package cmd

import (
	"log"

	"github.com/scrapnode/kanthor/project"
	"github.com/spf13/cobra"
)

func NewVersion() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "version of current release",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("kanthor %s", project.Version())
		},
	}

	return command
}
