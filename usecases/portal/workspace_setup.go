package portal

import (
	"context"
)

func (uc *workspace) Setup(ctx context.Context, req *WorkspaceSetupReq) (*WorkspaceSetupRes, error) {
	res, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		// starting with false
		status := map[string]bool{}

		for _, app := range req.Applications {
			status[app.Id] = false
		}
		for _, ep := range req.Endpoints {
			status[ep.Id] = false
		}
		for _, epr := range req.EndpointRules {
			status[epr.Id] = false
		}

		appIds, err := uc.repos.Application().BulkCreate(txctx, req.Applications)
		if err != nil {
			return nil, err
		}
		for _, appId := range appIds {
			status[appId] = true
		}

		epIds, err := uc.repos.Endpoint().BulkCreate(txctx, req.Endpoints)
		if err != nil {
			return nil, err
		}
		for _, epId := range epIds {
			status[epId] = true
		}

		eprIds, err := uc.repos.EndpointRule().BulkCreate(txctx, req.EndpointRules)
		if err != nil {
			return nil, err
		}
		for _, eprId := range eprIds {
			status[eprId] = true
		}

		res := &WorkspaceSetupRes{
			ApplicationIds:  appIds,
			EndpointIds:     epIds,
			EndpointRuleIds: eprIds,
			Status:          status,
		}
		return res, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*WorkspaceSetupRes), nil
}
