package administer

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
	"github.com/spf13/cobra"
)

func NewAttemptTrigger(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:  "attempt-trigger",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), isValidAppIdArg),
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := config.New(provider)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), time.Millisecond*time.Duration(conf.Trigger.Executor.Timeout))
			defer cancel()

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
			ds, err := datastore.New(provider)
			if err != nil {
				return err
			}

			defer func() {
				if err := db.Disconnect(ctx); err != nil {
					logger.Error(err)
				}
				if err := ds.Disconnect(ctx); err != nil {
					logger.Error(err)
				}
				if err := infra.Disconnect(ctx); err != nil {
					logger.Error(err)
				}
			}()

			if err := db.Connect(ctx); err != nil {
				return err
			}
			if err := ds.Connect(ctx); err != nil {
				return err
			}
			if err := infra.Connect(ctx); err != nil {
				return err
			}

			repos := repositories.New(logger, db, ds)
			uc := usecase.New(conf, logger, infra, repos)

			appId := args[0]
			app, err := repos.Database().Application().Get(ctx, appId)
			if err != nil {
				return err
			}

			// for selected message ids
			msgIds, err := cmd.Flags().GetStringArray("msg-id")
			if err != nil {
				return err
			}
			if len(msgIds) > 0 {
				return attemptTriggerWithMessageIds(cmd, ctx, conf, logger, uc, app, repos, msgIds)
			}

			return attemptTriggerWithDatetimeRange(cmd, ctx, conf, logger, infra, uc, app)
		},
	}

	f := time.Now().UTC().Format(format)
	command.Flags().StringP("from", "f", f, fmt.Sprintf("--from=%s (UTC +00:00) | beginning of scan time", f))

	t := time.Now().UTC().Add(time.Hour * 24).Format(format)
	command.Flags().StringP("to", "t", t, fmt.Sprintf("--to=%s (UTC +00:00) | end of scan time", t))

	command.Flags().IntP("concurrency", "", 500, "--concurrency=500 | concurrency exection items")
	command.Flags().StringArrayP("msg-id", "", []string{}, fmt.Sprintf("--msg-id=%s | select message to plan a trigger", entities.MsgId()))

	return command
}

func attemptTriggerWithMessageIds(
	cmd *cobra.Command,
	ctx context.Context,
	conf *config.Config,
	logger logging.Logger,
	uc usecase.Attempt,
	app *entities.ApplicationWithRelationship,
	repos repositories.Repositories,
	ids []string,
) error {
	messages, err := repos.Datastore().Message().ListByIds(ctx, app.Id, ids)
	if err != nil {
		return err
	}
	applicable, err := uc.Trigger().Applicable(ctx, app.Id)
	if err != nil {
		return err
	}

	out, err := uc.Trigger().Perform(ctx, app.Id, messages, applicable, conf.Trigger.Executor.ArrangeDelay)
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

	if len(out.Error) > 0 {
		for key, err := range out.Error {
			logger.Errorw(err.Error(), "key", key)
		}
	}

	return nil
}

func attemptTriggerWithDatetimeRange(
	cmd *cobra.Command,
	ctx context.Context,
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	uc usecase.Attempt,
	app *entities.ApplicationWithRelationship,
) error {
	concurrency, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return err
	}

	from, to, err := daterange(cmd)
	if err != nil {
		return err
	}

	in := &usecase.TriggerExecIn{
		Concurrency:  concurrency,
		ArrangeDelay: 0,
		Triggers: map[string]*entities.AttemptTrigger{
			suid.New("atttr"): {
				AppId: app.Id,
				Tier:  app.Workspace.Tier,
				From:  from.UnixMilli(),
				To:    to.UnixMilli(),
			},
		},
	}

	out, err := uc.Trigger().Exec(ctx, in)
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
}
