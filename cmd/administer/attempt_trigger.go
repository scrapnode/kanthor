package administer

import (
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
	"github.com/spf13/cobra"
)

var format = "2006-01-02"

func NewAttemptTrigger(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:  "attempt-trigger",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), isValidAppIdArg),
		RunE: func(cmd *cobra.Command, args []string) error {
			appId := args[0]

			from, to, err := daterange(cmd)
			if err != nil {
				return err
			}
			ctx, cancel, err := timeout(cmd)
			defer cancel()
			if err != nil {
				return err
			}

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
			app, err := repos.Database().Application().Get(ctx, appId)
			if err != nil {
				return err
			}

			uc := usecase.New(conf, logger, infra, repos)
			in := &usecase.TriggerExecIn{
				Concurrency:  100,
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
				"scheduled", len(out.Scheduled),
				"created", len(out.Created),
				"ko_count", len(out.Error),
			)

			return nil
		},
	}

	f := time.Now().UTC().Format(format)
	command.Flags().StringP("from", "f", f, fmt.Sprintf("--from=%s (UTC +00:00) | beginning of scan time", f))

	t := time.Now().UTC().Add(time.Hour * 24).Format(format)
	command.Flags().StringP("to", "t", t, fmt.Sprintf("--to=%s (UTC +00:00) | end of scan time", t))

	return command
}
