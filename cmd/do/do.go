package do

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "do",
	}

	command.AddCommand(NewImport(conf, logger))

	command.PersistentFlags().StringP("account-sub", "", "", "--account-sub=kanthor_root_key | select account sub to interact with")
	if err := command.MarkPersistentFlagRequired("account-sub"); err != nil {
		panic(err)
	}

	command.PersistentFlags().StringP("account-name", "", "", "--account-name=kanthor_root_key | specify name of the account you are using")

	return command
}
