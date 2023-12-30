package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointListMessageIn struct {
	*entities.ScanningQuery
	WsId string
	EpId string
}

func (in *EndpointListMessageIn) Validate() error {
	if err := in.ScanningQuery.Validate(); err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
	)
}

type EndpointListMessageOut struct {
	Data []entities.EndpointMessage
}

func (uc *endpoint) ListMessage(ctx context.Context, in *EndpointListMessageIn) (*EndpointListMessageOut, error) {
	ep, err := uc.repositories.Database().Endpoint().Get(ctx, in.WsId, in.EpId)
	if err != nil {
		return nil, err
	}

	reqMaps, err := uc.repositories.Datastore().Request().ListMessages(ctx, ep.Id, in.ScanningQuery)
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

	out := &EndpointListMessageOut{}
	for _, msg := range msgses {
		out.Data = append(out.Data, *uc.mapping(reqMaps, resMaps, msg))
	}
	return out, nil
}
