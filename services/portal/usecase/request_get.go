package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type RequestGetIn struct {
	WsId  string
	EpId  string
	MsgId string
	Id    string
}

func (in *RequestGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
		validator.StringStartsWith("msg_id", in.MsgId, entities.IdNsMsg),
		validator.StringStartsWith("id", in.Id, entities.IdNsReq),
	)
}

type RequestGetOut struct {
	Doc *entities.Request
}

func (uc *request) Get(ctx context.Context, in *RequestGetIn) (*RequestGetOut, error) {
	ep, err := uc.repositories.Database().Endpoint().Get(ctx, in.WsId, in.EpId)
	if err != nil {
		return nil, err
	}

	req, err := uc.repositories.Datastore().Request().Get(ctx, ep.Id, in.MsgId, in.Id)
	if err != nil {
		return nil, err
	}

	out := &RequestGetOut{Doc: req}
	return out, nil
}
