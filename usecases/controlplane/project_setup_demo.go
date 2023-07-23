package controlplane

import (
	"context"
)

func (uc *project) SetupDemo(ctx context.Context, req *ProjectSetupDemoReq) (*ProjectSetupDemoRes, error) {
	res, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ws, err := uc.repos.Workspace().Get(txctx, req.WorkspaceId)
		if err != nil {
			return nil, err
		}

		appIds, err := uc.repos.Application().BulkCreate(txctx, req.Applications)
		if err != nil {
			return nil, err
		}

		endpointIds, err := uc.repos.Endpoint().BulkCreate(txctx, req.Endpoints)
		if err != nil {
			return nil, err
		}

		endpointRuleIds, err := uc.repos.EndpointRule().BulkCreate(txctx, req.EndpointRules)
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
