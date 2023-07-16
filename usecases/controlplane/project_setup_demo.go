package controlplane

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/data/demo"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
)

func (usecase *project) SetupDemo(ctx context.Context, req *ProjectSetupDemoReq) (*ProjectSetupDemoRes, error) {
	ws, err := usecase.repos.Workspace().Get(ctx, req.WorkspaceId)
	if err != nil {
		return nil, err
	}

	owner := req.Account.Sub == ws.OwnerId
	if !owner {
		return nil, errors.New("only owner of this project can setup the demo")
	}

	// demo application
	app, err := usecase.repos.Application().Create(ctx, &entities.Application{
		WorkspaceId: ws.Id,
		Name:        constants.DemoApplicationName,
	})
	if err != nil {
		return nil, err
	}

	// demo endpoints
	endpointIds, err := usecase.repos.Endpoint().BulkCreate(ctx, demo.Endpoints(app.Id))
	if err != nil {
		return nil, err
	}

	// demo rules for endpoints
	endpointRuleIds, err := usecase.repos.EndpointRule().BulkCreate(ctx, demo.EndpointRules(app.Id, endpointIds))
	if err != nil {
		return nil, err
	}

	// must clear the cache because of new endpoints and rules
	cacheKey := cache.Key("APP_WITH_ENDPOINTS", app.Id)
	if err := usecase.cache.Del(cacheKey); err != nil {
		return nil, err
	}

	res := &ProjectSetupDemoRes{ApplicationId: app.Id, EndpointIds: endpointIds, EndpointRuleIds: endpointRuleIds}
	return res, nil
}
