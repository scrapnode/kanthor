package get

import (
	"context"

	"github.com/scrapnode/kanthor/cmd/client/utils"
	"github.com/spf13/cobra"
)

func Account() *cobra.Command {
	command := &cobra.Command{
		Use: "account",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := utils.ClientFromCmd(cmd)
			if err != nil {
				return err
			}
			res, err := client.Account.Get(context.Background())
			if err != nil {
				return err
			}

			return utils.ResponseTo(cmd, res)
		},
	}

	return command
}
