package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleListIn struct {
	*entities.Query
	WsId  string
	AppId string
	EpId  string
}

func (in *EndpointRuleListIn) Validate() error {
	if err := in.Query.Validate(); err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWithIfNotEmpty("app_id", in.AppId, entities.IdNsApp),
		validator.StringStartsWithIfNotEmpty("ep_id", in.EpId, entities.IdNsEp),
	)
}

type EndpointRuleListOut struct {
	Data  []entities.EndpointRule
	Count int64
}

func (uc *endpointRule) List(ctx context.Context, in *EndpointRuleListIn) (*EndpointRuleListOut, error) {
	data, err := uc.repositories.EndpointRule().List(ctx, in.WsId, in.AppId, in.EpId, in.Query)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.EndpointRule().Count(ctx, in.WsId, in.AppId, in.EpId, in.Query)
	if err != nil {
		return nil, err
	}

	out := &EndpointRuleListOut{Data: data, Count: count}
	return out, nil
}
