package account

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/services/ioc"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/spf13/cobra"
	"time"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:  "account",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			// hoisting flags
			file, err := cmd.Flags().GetString("data")
			if err != nil {
				return err
			}

			withCreds, err := cmd.Flags().GetBool("generate-credentials")
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), time.Minute*2)
			defer cancel()

			meter, err := metric.NewNoop(nil, logger)
			if err != nil {
				return err
			}
			uc, err := ioc.InitializePortalUsecase(conf, logger, meter)
			if err != nil {
				return err
			}
			if err := uc.Connect(ctx); err != nil {
				return err
			}

			coord, err := coordinator.New(&conf.Coordinator, logger)
			if err != nil {
				return err
			}
			if err := coord.Connect(ctx); err != nil {
				return err
			}

			defer func() {
				if err := uc.Disconnect(ctx); err != nil {
					logger.Error(err)
				}

				if err := coord.Disconnect(ctx); err != nil {
					logger.Error(err)
				}

			}()

			account, err := uc.Account().Setup(ctx, &usecase.AccountSetupReq{AccountId: args[0]})
			if err != nil {
				return err
			}

			if err := apps(uc, ctx, account.Workspace, file); err != nil {
				return err
			}

			if err := creds(uc, ctx, account.Workspace, withCreds); err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().StringP("data", "", "", "--data=./data/interchange/demo.json | workspace data of setup account")
	command.Flags().BoolP("generate-credentials", "", false, "--generate-credentials | generate new credentials for the workspace of setup account")

	return command
}
