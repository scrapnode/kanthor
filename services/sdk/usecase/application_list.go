package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationListIn struct {
	*entities.PagingQuery
	WsId string
}

func (in *ApplicationListIn) Validate() error {
	if err := in.PagingQuery.Validate(); err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
	)
}

type ApplicationListOut struct {
	Data  []entities.Application
	Count int64
}

func (uc *application) List(ctx context.Context, in *ApplicationListIn) (*ApplicationListOut, error) {
	data, err := uc.repositories.Application().List(ctx, in.WsId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.Application().Count(ctx, in.WsId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	out := &ApplicationListOut{Data: data, Count: count}
	return out, nil
}
