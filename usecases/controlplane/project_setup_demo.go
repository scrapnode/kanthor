package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (usecase *project) SetupDemo(ctx context.Context, req *ProjectSetupDemoReq) (*ProjectSetupDemoRes, error) {
	res, err := usecase.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ws, err := usecase.repos.Workspace().Get(txctx, req.WorkspaceId)
		if err != nil {
			return nil, err
		}

		// demo applications
		var apps []entities.Application
		for _, app := range req.Entities.Applications {
			app.GenId()
			app.WorkspaceId = ws.Id
			app.SetAT(usecase.timer.Now())
			apps = append(apps, app)
		}
		appIds, err := usecase.repos.Application().BulkCreate(txctx, apps)
		if err != nil {
			return nil, err
		}

		// demo endpoints
		var endpoints []entities.Endpoint
		for _, endpoint := range req.Entities.Endpoints {
			endpoint.GenId()
			endpoint.SetAT(usecase.timer.Now())
			endpoints = append(endpoints, endpoint)
		}
		endpointIds, err := usecase.repos.Endpoint().BulkCreate(txctx, endpoints)
		if err != nil {
			return nil, err
		}

		// demo rules for endpoints
		var rules []entities.EndpointRule
		for _, rule := range req.Entities.EndpointRules {
			rule.GenId()
			rule.SetAT(usecase.timer.Now())
			rules = append(rules, rule)
		}
		endpointRuleIds, err := usecase.repos.EndpointRule().BulkCreate(txctx, rules)
		if err != nil {
			return nil, err
		}

		res := &ProjectSetupDemoRes{ApplicationIds: appIds, EndpointIds: endpointIds, EndpointRuleIds: endpointRuleIds}
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*ProjectSetupDemoRes), nil
}
