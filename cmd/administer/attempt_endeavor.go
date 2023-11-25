package administer

import (
	"fmt"
	"log"
	"time"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/database"
	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/attempt/config"
	"github.com/scrapnode/kanthor/services/attempt/repositories"
	"github.com/scrapnode/kanthor/services/attempt/usecase"
	"github.com/spf13/cobra"
)

func NewAttemptEndeavor(provider configuration.Provider) *cobra.Command {
	command := &cobra.Command{
		Use:  "attempt-endeavor",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, to, err := daterange(cmd)
			if err != nil {
				return err
			}
			ctx, cancel, err := timeout(cmd)
			defer cancel()
			if err != nil {
				return err
			}
			concurrency, err := cmd.Flags().GetInt("concurrency")
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
			uc := usecase.New(conf, logger, infra, repos)

			ch := repos.Datastore().Attempt().Scan(ctx, *from, *to, time.Now().UTC().UnixMilli(), concurrency)
			for r := range ch {
				if r.Error != nil {
					return r.Error
				}

				strive, err := uc.Endeavor().Evaluate(ctx, r.Data)
				if err != nil {
					return err
				}
				if err := repos.Datastore().Attempt().MarkIgnore(ctx, strive.Ignore); err != nil {
					logger.Errorw("unable to ignore attempts", "req_ids", strive.Ignore)
				}

				in := &usecase.EndeavorExecIn{
					Concurrency: concurrency,
					Attempts:    map[string]*entities.Attempt{},
				}

				// The iteration variables may be declared by the "range" clause using a form of short variable declaration (:=).
				// In this case their types are set to the types of the respective iteration values and their scope is the block of the "for" statement;
				// **they are re-used in each iteration**. If the iteration variables are declared outside the "for" statement,
				// after execution their values will be those of the last iteration.
				for i := 0; i < len(strive.Attemptable); i++ {
					in.Attempts[strive.Attemptable[i].ReqId] = &strive.Attemptable[i]
				}

				out, err := uc.Endeavor().Exec(ctx, in)
				if err != nil {
					return err
				}

				log.Printf("out.Success -> %d", len(out.Success))
				logger.Infow(
					"administer attempt trigger",
					"from", from.Format(time.RFC3339),
					"to", to.Format(time.RFC3339),
					"ok_count", len(out.Success),
					"rescheduled_count", len(out.Rescheduled),
					"completed_count", len(out.Completed),
					"ko_count", len(out.Error),
				)

			}

			return nil
		},
	}

	f := time.Now().UTC().Format(format)
	command.Flags().StringP("from", "f", f, fmt.Sprintf("--from=%s (UTC +00:00) | beginning of scan time", f))

	t := time.Now().UTC().Add(time.Hour * 24).Format(format)
	command.Flags().StringP("to", "t", t, fmt.Sprintf("--to=%s (UTC +00:00) | end of scan time", t))

	command.Flags().IntP("concurrency", "", 100, "--concurrency=500 | concurrency exection items")

	return command
}