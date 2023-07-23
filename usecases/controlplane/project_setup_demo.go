package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *project) SetupDemo(ctx context.Context, req *ProjectSetupDemoReq) (*ProjectSetupDemoRes, error) {
	res, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ws, err := uc.repos.Workspace().Get(txctx, req.WorkspaceId)
		if err != nil {
			return nil, err
		}

		// demo applications
		apps := []entities.Application{}
		for _, app := range req.Applications {
			app.GenId()
			app.WorkspaceId = ws.Id
			app.SetAT(uc.timer.Now())
			apps = append(apps, app)
		}
		appIds, err := uc.repos.Application().BulkCreate(txctx, apps)
		if err != nil {
			return nil, err
		}

		// demo endpoints
		endpoints := []entities.Endpoint{}
		for _, endpoint := range req.Endpoints {
			endpoint.GenId()
			endpoint.SetAT(uc.timer.Now())
			endpoints = append(endpoints, endpoint)
		}
		endpointIds, err := uc.repos.Endpoint().BulkCreate(txctx, endpoints)
		if err != nil {
			return nil, err
		}

		// demo rules for endpoints
		rules := []entities.EndpointRule{}
		for _, rule := range req.EndpointRules {
			rule.GenId()
			rule.SetAT(uc.timer.Now())
			rules = append(rules, rule)
		}
		endpointRuleIds, err := uc.repos.EndpointRule().BulkCreate(txctx, rules)
		if err != nil {
			return nil, err
		}

		res := &ProjectSetupDemoRes{
			WorkspaceId:     ws.Id,
			WorkspaceTier:   ws.Tier.Name,
			ApplicationIds:  appIds,
			EndpointIds:     endpointIds,
			EndpointRuleIds: endpointRuleIds,
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*ProjectSetupDemoRes), nil
}
