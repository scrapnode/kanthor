package cmd

import (
	"github.com/scrapnode/kanthor/cmd/show"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/spf13/cobra"
)

func NewVersion(provider configuration.Provider, conf *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "version of current release, shortcut of show.version",
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			return show.Version(conf, verbose)
		},
	}

	return command
}
