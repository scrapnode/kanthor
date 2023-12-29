package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type RequestListIn struct {
	*entities.ScanningQuery
	WsId string
	EpId string
}

func (in *RequestListIn) Validate() error {
	if err := in.ScanningQuery.Validate(); err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
	)
}

type RequestListOut struct {
	Data []entities.Request
}

func (uc *request) List(ctx context.Context, in *RequestListIn) (*RequestListOut, error) {
	ep, err := uc.repositories.Database().Endpoint().Get(ctx, in.WsId, in.EpId)
	if err != nil {
		return nil, err
	}

	data, err := uc.repositories.Datastore().Request().List(ctx, ep.Id, in.ScanningQuery)
	if err != nil {
		return nil, err
	}

	out := &RequestListOut{Data: data}
	return out, nil
}
