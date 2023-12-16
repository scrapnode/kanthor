package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type TriggerExecWithDateRangeIn struct {
	AppId        string
	ArrangeDelay int64
	Concurrency  int

	From int64
	To   int64
}

func (in *TriggerExecWithDateRangeIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.NumberGreaterThan("concurrency", in.Concurrency, 0),
		validator.NumberGreaterThan("arrange_delay", in.ArrangeDelay, 60000),
		validator.NumberLessThan("from", in.From, in.To),
		validator.NumberLessThanOrEqual("from", in.To, time.Now().UTC().UnixMilli()),
	)
}

func (uc *cli) TriggerExecWithDateRange(ctx context.Context, in *TriggerExecWithDateRangeIn) (*TriggerExecOut, error) {
	app, err := uc.repositories.Database().Application().Get(ctx, in.AppId)
	if err != nil {
		return nil, err
	}

	i := &TriggerExecIn{
		Concurrency:  in.Concurrency,
		ArrangeDelay: in.ArrangeDelay,
		Triggers: map[string]*entities.AttemptTrigger{
			suid.New("atttr"): {
				AppId: app.Id,
				Tier:  app.Workspace.Tier,
				From:  in.From,
				To:    in.To,
			},
		},
	}
	if err := i.Validate(); err != nil {
		return nil, err
	}

	return uc.trigger.Exec(ctx, i)
}
