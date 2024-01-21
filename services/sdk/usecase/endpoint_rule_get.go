package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleGetIn struct {
	WsId string
	Id   string
}

func (in *EndpointRuleGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsEpr),
	)
}

type EndpointRuleGetOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Get(ctx context.Context, in *EndpointRuleGetIn) (*EndpointRuleGetOut, error) {
	epr, err := uc.repositories.Database().EndpointRule().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}

	return &EndpointRuleGetOut{Doc: epr}, nil
}
