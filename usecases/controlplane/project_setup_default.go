package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (usecase *project) SetupDefault(ctx context.Context, req *ProjectSetupDefaultReq) (*ProjectSetupDefaultRes, error) {
	entity := &entities.Workspace{
		OwnerId: req.Account.Sub,
		Name:    req.WorkspaceName,
	}
	if entity.Name == "" {
		entity.Name = constants.DefaultWorkspaceName
	}
	entity.Tier = &entities.WorkspaceTier{Name: req.WorkspaceTier}
	if entity.Tier.Name == "" {
		entity.Tier.Name = constants.DefaultWorkspaceTier
	}

	ws, err := usecase.repos.Workspace().Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	res := &ProjectSetupDefaultRes{WorkspaceId: ws.Id, WorkspaceTier: ws.Tier.Name}
	return res, nil
}
