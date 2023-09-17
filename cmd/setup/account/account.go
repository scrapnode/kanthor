package account

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/validation"
	"github.com/scrapnode/kanthor/services/ioc"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
	"github.com/spf13/cobra"
)

func New(conf *config.Config, logger logging.Logger) *cobra.Command {
	validator := validation.New()

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
			dest, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), time.Minute*2)
			defer cancel()

			meter, err := metric.NewNoop(nil, logger)
			if err != nil {
				return err
			}
			uc, err := ioc.InitializePortalUsecase(conf, logger, validator, meter)
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

			out := &output{json: map[string]any{}}

			account, err := uc.Account().Setup(ctx, &usecase.AccountSetupReq{AccountId: args[0]})
			if err != nil {
				return err
			}
			out.AddJson("id", account.Workspace.Id)
			out.AddJson("tier", account.WorkspaceTier.Name)

			if err := apps(validator, uc, ctx, account.Workspace, file, out); err != nil {
				return err
			}

			if withCreds {
				if err := creds(validator, coord, uc, ctx, account.Workspace, out); err != nil {
					return err
				}
			}

			return out.Render(dest)
		},
	}

	command.Flags().StringP("data", "", "", "--data=./data/interchange/demo.json | workspace data of setup account")
	command.Flags().StringP("output", "o", "", "--out=./output.json | either json file or stdout (if no file path is set)")
	command.Flags().BoolP("generate-credentials", "", false, "--generate-credentials | generate new credentials for the workspace of setup account")

	return command
}
