package usecase

import (
	"context"
	"errors"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services/portal/repositories/ds"
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
	Doc       *entities.EndpointMessage
	Requests  []entities.Request
	Responses []entities.Response
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

	resMaps, err := uc.repositories.Datastore().Response().GetMessages(ctx, ep.Id, reqMaps.MsgIds)
	if err != nil {
		return nil, err
	}

	msgses, err := uc.repositories.Datastore().Message().GetByIds(ctx, ep.AppId, reqMaps.MsgIds)
	if err != nil {
		return nil, err
	}

	if len(msgses) == 0 {
		return nil, errors.New("message was not found")
	}

	msg := msgses[0]
	out := &EndpointGetMessageOut{
		Doc:       uc.mapping(reqMaps, resMaps, msg),
		Requests:  reqMaps.Maps[msgses[0].Id],
		Responses: resMaps.Maps[msgses[0].Id],
	}
	if _, has := reqMaps.Maps[msg.Id]; has {
		out.Requests = reqMaps.Maps[msg.Id]
	}
	if _, has := resMaps.Maps[msg.Id]; has {
		out.Responses = resMaps.Maps[msg.Id]
	}
	return out, nil
}

func (uc *endpoint) mapping(reqMaps *ds.MessageRequestMaps, resMaps *ds.MessageResponsetMaps, msg entities.Message) *entities.EndpointMessage {
	data := &entities.EndpointMessage{Message: msg}

	if _, has := reqMaps.Maps[msg.Id]; has {
		data.RequestCount = len(reqMaps.Maps[msg.Id])
		if data.RequestCount > 0 {
			data.RequestLatestTs = reqMaps.Maps[msg.Id][0].Timestamp
		}
	}

	if _, has := resMaps.Maps[msg.Id]; has {
		data.ResponseCount = len(resMaps.Maps[msg.Id])
		if data.ResponseCount > 0 {
			data.ResponseLatestTs = resMaps.Maps[msg.Id][0].Timestamp
		}
	}

	if id, has := resMaps.Success[msg.Id]; has {
		data.SuccessId = id
	}

	return data
}
