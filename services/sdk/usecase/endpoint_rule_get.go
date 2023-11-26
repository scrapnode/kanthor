package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleGetIn struct {
	EpId string
	Id   string
}

func (req *EndpointRuleGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ep_id", req.EpId, entities.IdNsEp),
		validator.StringStartsWith("id", req.EpId, entities.IdNsEpr),
	)
}

type EndpointRuleGetOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Get(ctx context.Context, req *EndpointRuleGetIn) (*EndpointRuleGetOut, error) {
	ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

	key := CacheKeyEpr(req.EpId, req.Id)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*EndpointRuleGetOut, error) {
		ep, err := uc.repositories.Endpoint().GetOfWorkspace(ctx, ws, req.EpId)
		if err != nil {
			return nil, err
		}

		epr, err := uc.repositories.EndpointRule().Get(ctx, ep, req.Id)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}
		return &EndpointRuleGetOut{Doc: epr}, nil
	})
}
