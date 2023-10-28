package account

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/portal/config"
	"github.com/scrapnode/kanthor/services/portal/repositories"
	"github.com/scrapnode/kanthor/services/portal/usecase"
	"github.com/spf13/cobra"
)

func New(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:  "account",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			// hoisting flags
			file, err := cmd.Flags().GetString("data")
			if err != nil {
				return err
			}
			dest, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), time.Minute*5)
			defer cancel()

			conf, err := config.New(provider)
			if err != nil {
				return err
			}
			logger, err := logging.New(provider)
			if err != nil {
				return err
			}
			infra, err := infrastructure.New(provider)
			if err != nil {
				return err
			}
			db, err := database.New(provider)
			if err != nil {
				return err
			}

			defer func() {
				if err := db.Disconnect(ctx); err != nil {
					logger.Error(err)
				}
				if err := infra.Disconnect(ctx); err != nil {
					logger.Error(err)
				}
			}()

			if err := db.Connect(ctx); err != nil {
				return err
			}
			if err := infra.Connect(ctx); err != nil {
				return err
			}

			repos := repositories.New(logger, db)
			uc := usecase.New(conf, logger, infra, repos)

			out := &output{json: map[string]any{}}

			account, err := uc.Account().Setup(ctx, &usecase.AccountSetupReq{AccountId: args[0]})
			if err != nil {
				return err
			}
			out.AddJson("id", account.Workspace.Id)
			out.AddJson("tier", account.Workspace.Tier)

			if err := apps(uc, ctx, account.Workspace, file, out); err != nil {
				return err
			}
			if err := creds(uc, ctx, account.Workspace, out); err != nil {
				return err
			}

			return out.Render(dest)
		},
	}

	command.Flags().StringP("data", "", "", "--data=./data/interchange/demo.json | workspace data of setup account")
	command.Flags().StringP("output", "o", "", "--out=./output.json | either json file or stdout (if no file path is set)")

	return command
}
