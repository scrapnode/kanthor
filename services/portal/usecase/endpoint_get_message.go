package usecase

import (
	"context"
	"errors"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointGetMessageIn struct {
	WsId  string
	EpId  string
	MsgId string
}

func (in *EndpointGetMessageIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
		validator.StringStartsWith("msg_id", in.MsgId, entities.IdNsMsg),
	)
}

type EndpointGetMessageOut struct {
	Doc *entities.EndpointMessage
}

func (uc *endpoint) GetMessage(ctx context.Context, in *EndpointGetMessageIn) (*EndpointGetMessageOut, error) {
	ep, err := uc.repositories.Database().Endpoint().Get(ctx, in.WsId, in.EpId)
	if err != nil {
		return nil, err
	}

	reqMaps, err := uc.repositories.Datastore().Request().GetMessage(ctx, ep.Id, in.MsgId)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	resMaps, err := uc.repositories.Datastore().Response().ListMessages(ctx, ep.Id, reqMaps.MsgIds)
	if err != nil {
		return nil, err
	}

	msgses, err := uc.repositories.Datastore().Message().ListByIds(ctx, ep.AppId, reqMaps.MsgIds)
	if err != nil {
		return nil, err
	}

	if len(msgses) == 0 {
		return nil, errors.New("message was not found")
	}

	out := &EndpointGetMessageOut{Doc: uc.mapping(reqMaps, resMaps, msgses[0])}
	return out, nil
}
