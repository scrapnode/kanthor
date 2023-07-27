package portal

import "context"

func (uc *workspace) Import(ctx context.Context, req *WorkspaceImportReq) (*WorkspaceImportRes, error) {
	res, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		wsIds, err := uc.repos.Workspace().BulkCreate(txctx, req.Workspaces)
		if err != nil {
			return nil, err
		}
		wstIds, err := uc.repos.WorkspaceTier().BulkCreate(txctx, req.WorkspaceTiers)
		if err != nil {
			return nil, err
		}
		wscIds, err := uc.repos.WorkspaceCredentials().BulkCreate(txctx, req.WorkspaceCredentials)
		if err != nil {
			return nil, err
		}
		appIds, err := uc.repos.Application().BulkCreate(txctx, req.Applications)
		if err != nil {
			return nil, err
		}
		epIds, err := uc.repos.Endpoint().BulkCreate(txctx, req.Endpoints)
		if err != nil {
			return nil, err
		}
		eprIds, err := uc.repos.EndpointRule().BulkCreate(txctx, req.EndpointRules)
		if err != nil {
			return nil, err
		}

		res := &WorkspaceImportRes{
			WorkspaceIds:            wsIds,
			WorkspaceTierIds:        wstIds,
			WorkspaceCredentialsIds: wscIds,
			ApplicationIds:          appIds,
			EndpointIds:             epIds,
			EndpointRuleIds:         eprIds,
		}
		return res, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*WorkspaceImportRes), nil
}
