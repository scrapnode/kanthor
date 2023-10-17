package sdk

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointRuleGetReq struct {
	EpId string
	Id   string
}

func (req *EndpointRuleGetReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ep_id", req.EpId, entities.IdNsEp),
		validator.StringStartsWith("id", req.EpId, entities.IdNsEpr),
	)
}

type EndpointRuleGetRes struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Get(ctx context.Context, req *EndpointRuleGetReq) (*EndpointRuleGetRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	key := CacheKeyEpr(req.EpId, req.Id)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*EndpointRuleGetRes, error) {
		ep, err := uc.repos.Endpoint().GetOfWorkspace(ctx, ws, req.EpId)
		if err != nil {
			return nil, err
		}

		epr, err := uc.repos.EndpointRule().Get(ctx, ep, req.Id)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}
		res := &EndpointRuleGetRes{Doc: epr}
		return res, nil
	})
}
