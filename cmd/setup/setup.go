package setup

import (
	"errors"
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:       "setup",
		ValidArgs: []string{"demo"},
		Args:      cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}

			sub, err := cmd.Flags().GetString("account-sub")
			if err != nil {
				return err
			}

			name := args[0]
			if name == "demo" {
				input, err := cmd.Flags().GetString("demo-input-file")
				if err != nil {
					return err
				}
				if input == "" {
					return errors.New("setup.demo: demo input file must be provided")
				}
				return Demo(conf, logger, sub, input, verbose)
			}

			return fmt.Errorf("setup: unknow setup [%s]", name)
		},
	}

	command.Flags().StringP("account-sub", "", "", "--account-sub=kanthor_root_key | select account to setup stuffs")
	if err := command.MarkFlagRequired("account-sub"); err != nil {
		panic(err)
	}

	// demo specific flags
	command.Flags().StringP("demo-input-file", "", "", "--demo-input-file=data/demo/project.json | the json input file to create demo project")

	return command
}
