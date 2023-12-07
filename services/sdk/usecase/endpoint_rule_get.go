package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleGetIn struct {
	Ws   *entities.Workspace
	EpId string
	Id   string
}

func (in *EndpointRuleGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.PointerNotNil("ws", in.Ws),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
		validator.StringStartsWith("id", in.EpId, entities.IdNsEpr),
	)
}

type EndpointRuleGetOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Get(ctx context.Context, in *EndpointRuleGetIn) (*EndpointRuleGetOut, error) {
	key := CacheKeyEpr(in.EpId, in.Id)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*EndpointRuleGetOut, error) {
		ep, err := uc.repositories.Endpoint().GetOfWorkspace(ctx, in.Ws, in.EpId)
		if err != nil {
			return nil, err
		}

		epr, err := uc.repositories.EndpointRule().Get(ctx, ep, in.Id)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}
		return &EndpointRuleGetOut{Doc: epr}, nil
	})
}
