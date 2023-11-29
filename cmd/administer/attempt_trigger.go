package administer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scrapnode/kanthor/cmd/utils"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
	"github.com/spf13/cobra"
)

func NewAttemptTrigger(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:  "attempt-trigger",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), isValidAppIdArg),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := logging.New(provider)
			if err != nil {
				return err
			}
			conf, err := config.New(provider)
			if err != nil {
				return err
			}

			cli, err := CommandLine(provider, services.ATTEMPT_TRIGGER_CLI)
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			if err = cli.Start(ctx); err != nil {
				return err
			}

			defer func() {
				err = utils.Stop(cli)
			}()

			appId := args[0]
			uc := cli.Usecase().(usecase.Cli)

			ids, err := cmd.Flags().GetStringSlice("msg-id")
			if err != nil {
				return err
			}

			if len(ids) > 0 {
				in := &usecase.TriggerExecWithMessageIdsIn{
					AppId:        appId,
					ArrangeDelay: conf.Trigger.Executor.ArrangeDelay,
					MsgIds:       ids,
				}

				out, err := uc.TriggerExecWithMessageIds(ctx, in)
				if err != nil {
					return err
				}

				logger.Infow(
					"administer attempt trigger",
					"ids", ids,
					"ok_count", len(out.Success),
					"scheduled_count", len(out.Scheduled),
					"created_count", len(out.Created),
					"ko_count", len(out.Error),
				)
				return nil
			}

			from, to, err := utils.DateRange(cmd)
			if err != nil {
				return err
			}

			in := &usecase.TriggerExecWithDateRangeIn{
				AppId:        appId,
				ArrangeDelay: conf.Trigger.Executor.ArrangeDelay,
				Concurrency:  conf.Trigger.Executor.Concurrency,
				From:         from.UnixMilli(),
				To:           to.UnixMilli(),
			}
			out, err := uc.TriggerExecWithDateRange(ctx, in)
			if err != nil {
				return err
			}
			logger.Infow(
				"administer attempt trigger",
				"from", from.Format(time.RFC3339),
				"to", to.Format(time.RFC3339),
				"ok_count", len(out.Success),
				"scheduled_count", len(out.Scheduled),
				"created_count", len(out.Created),
				"ko_count", len(out.Error),
			)
			return nil
		},
	}

	f := time.Now().UTC().Format(format)
	command.Flags().StringP("from", "f", f, fmt.Sprintf("--from=%s (UTC +00:00) | beginning of scan time", f))

	t := time.Now().UTC().Add(time.Hour * 24).Format(format)
	command.Flags().StringP("to", "t", t, fmt.Sprintf("--to=%s (UTC +00:00) | end of scan time", t))

	command.Flags().StringSliceP("msg-id", "", []string{}, fmt.Sprintf("--msg-id=%s | select message to plan", suid.New(entities.IdNsMsg)))

	return command
}
