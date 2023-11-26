package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/assessor"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
)

func (uc *trigger) Applicable(ctx context.Context, appId string) (*assessor.Assets, error) {
	key := utils.Key("scheduler", appId)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour, func() (*assessor.Assets, error) {
		endpoints, err := uc.repositories.Database().Endpoint().List(ctx, appId)
		if err != nil {
			return nil, err
		}
		returning := &assessor.Assets{EndpointMap: map[string]entities.Endpoint{}}
		for _, ep := range endpoints {
			returning.EndpointMap[ep.Id] = ep
		}

		rules, err := uc.repositories.Database().Endpoint().Rules(ctx, appId)
		if err != nil {
			return nil, err
		}
		returning.Rules = rules

		return returning, nil
	})
}
