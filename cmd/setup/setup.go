package setup

import (
	"github.com/scrapnode/kanthor/cmd/setup/account"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use: "setup",
	}

	command.AddCommand(account.New(conf, logger))

	return command
}
