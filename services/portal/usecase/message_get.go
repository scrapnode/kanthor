package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type MessageGetIn struct {
	WsId  string
	AppId string
	Id    string
}

func (in *MessageGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.StringStartsWith("id", in.Id, entities.IdNsMsg),
	)
}

type MessageGetOut struct {
	Doc *entities.Message
}

func (uc *message) Get(ctx context.Context, in *MessageGetIn) (*MessageGetOut, error) {
	app, err := uc.repositories.Database().Application().Get(ctx, in.WsId, in.AppId)
	if err != nil {
		return nil, err
	}

	msg, err := uc.repositories.Datastore().Message().Get(ctx, app.Id, in.Id)
	if err != nil {
		return nil, err
	}

	out := &MessageGetOut{Doc: msg}
	return out, nil
}
