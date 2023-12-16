package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type TriggerExecWithMessageIdsIn struct {
	AppId        string
	ArrangeDelay int64

	MsgIds []string
}

func (in *TriggerExecWithMessageIdsIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.NumberGreaterThan("arrange_delay", in.ArrangeDelay, 60000),
		validator.SliceRequired[string]("msg_ids", in.MsgIds),
		validator.Slice[string](in.MsgIds, func(_ int, id *string) error {
			return validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp)()
		}),
	)
}

func (uc *cli) TriggerExecWithMessageIds(ctx context.Context, in *TriggerExecWithMessageIdsIn) (*TriggerExecOut, error) {
	app, err := uc.repositories.Database().Application().Get(ctx, in.AppId)
	if err != nil {
		return nil, err
	}

	applicable, err := uc.trigger.Applicable(ctx, app.Id)
	if err != nil {
		return nil, err
	}

	messages, err := uc.repositories.Datastore().Message().ListByIds(ctx, app.Id, in.MsgIds)
	if err != nil {
		return nil, err
	}

	i := &TriggerPerformIn{
		AppId:        app.Id,
		Concurrency:  len(messages),
		ArrangeDelay: in.ArrangeDelay,
		Applicable:   applicable,
		Messages:     messages,
	}
	if err := i.Validate(); err != nil {
		return nil, err
	}
	return uc.trigger.Perform(ctx, i)
}
