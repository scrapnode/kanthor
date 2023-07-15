package show

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider, conf *config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:       "show",
		ValidArgs: []string{"config", "version"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}
			validating, err := cmd.Flags().GetBool("validate")
			if err != nil {
				return err
			}

			name := args[0]
			if name == "config" {
				return Config(conf, provider.Sources(), verbose, validating)
			}
			if name == "version" {
				return Version(conf, verbose)
			}

			return nil
		},
	}

	command.Flags().BoolP("validate", "", false, "should validate the output we show you")
	return command
}
