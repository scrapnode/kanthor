package setup

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "setup",
	}

	command.AddCommand(NewDemo(conf, logger))

	command.PersistentFlags().StringP("account-sub", "", "", "--account-sub=kanthor_root_key | select account to setup stuffs")
	if err := command.MarkPersistentFlagRequired("account-sub"); err != nil {
		panic(err)
	}
	return command
}
