package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointListIn struct {
	*entities.Query
	WsId  string
	AppId string
}

func (in *EndpointListIn) Validate() error {
	if err := in.Query.Validate(); err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWithIfNotEmpty("app_id", in.AppId, entities.IdNsApp),
	)
}

type EndpointListOut struct {
	Data  []entities.Endpoint
	Count int64
}

func (uc *endpoint) List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error) {
	data, err := uc.repositories.Endpoint().List(ctx, in.WsId, in.AppId, in.Query)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.Endpoint().Count(ctx, in.WsId, in.AppId, in.Query)
	if err != nil {
		return nil, err
	}

	out := &EndpointListOut{Data: data, Count: count}
	return out, nil
}
