package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type AnalyticsGetOverviewIn struct {
	WsId string
}

func (in *AnalyticsGetOverviewIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
	)
}

type AnalyticsGetOverviewOut struct {
	CredentialsCount int64
	ApplicationCount int64
	EndpointCount    int64
	RuleCount        int64
}

func (uc *analytics) GetOverview(ctx context.Context, in *AnalyticsGetOverviewIn) (*AnalyticsGetOverviewOut, error) {
	credCount, err := uc.repositories.Database().WorkspaceCredentials().Count(ctx, in.WsId, entities.DefaultPagingQuery)
	if err != nil {
		return nil, err
	}
	appCount, err := uc.repositories.Database().Application().Count(ctx, in.WsId, entities.DefaultPagingQuery)
	if err != nil {
		return nil, err
	}
	epCount, err := uc.repositories.Database().Endpoint().Count(ctx, in.WsId, entities.DefaultPagingQuery)
	if err != nil {
		return nil, err
	}

	eprCount, err := uc.repositories.Database().EndpointRule().Count(ctx, in.WsId, entities.DefaultPagingQuery)
	if err != nil {
		return nil, err
	}

	out := &AnalyticsGetOverviewOut{
		CredentialsCount: credCount,
		ApplicationCount: appCount,
		EndpointCount:    epCount,
		RuleCount:        eprCount,
	}
	return out, nil
}
