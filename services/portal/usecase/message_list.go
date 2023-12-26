package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type MessageListIn struct {
	*entities.ScanningQuery
	WsId  string
	AppId string
}

func (in *MessageListIn) Validate() error {
	if err := in.ScanningQuery.Validate(); err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
	)
}

type MessageListOut struct {
	Data []entities.Message
}

func (uc *message) List(ctx context.Context, in *MessageListIn) (*MessageListOut, error) {
	app, err := uc.repositories.Database().Application().Get(ctx, in.WsId, in.AppId)
	if err != nil {
		return nil, err
	}

	data, err := uc.repositories.Datastore().Message().List(ctx, app.Id, in.ScanningQuery)
	if err != nil {
		return nil, err
	}

	out := &MessageListOut{Data: data}
	return out, nil
}