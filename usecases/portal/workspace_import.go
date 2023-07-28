package portal

import (
	"context"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *workspace) Import(ctx context.Context, req *WorkspaceImportReq) (*WorkspaceImportRes, error) {
	// by default all imported workspace must be start with default tier
	// the reason why we have to do that is because of security risk
	// let image a customer export data with their workspace tier
	// they can change it to whatever tier they want then import back to the system
	now := uc.timer.Now()
	var tiers []entities.WorkspaceTier
	for _, ws := range req.Workspaces {
		tier := entities.WorkspaceTier{WorkspaceId: ws.Id, Name: constants.DefaultWorkspaceTier}
		tier.GenId()
		tier.SetAT(now)
		tiers = append(tiers, tier)
	}

	res, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		wsIds, err := uc.repos.Workspace().BulkCreate(txctx, req.Workspaces)
		if err != nil {
			return nil, err
		}
		wstIds, err := uc.repos.WorkspaceTier().BulkCreate(txctx, tiers)
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
			WorkspaceIds:     wsIds,
			WorkspaceTierIds: wstIds,
			ApplicationIds:   appIds,
			EndpointIds:      epIds,
			EndpointRuleIds:  eprIds,
		}
		return res, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*WorkspaceImportRes), nil
}
