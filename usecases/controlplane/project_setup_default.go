package controlplane

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/database"
)

func (uc *project) SetupDefault(ctx context.Context, req *ProjectSetupDefaultReq) (*ProjectSetupDefaultRes, error) {
	res, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		existing, err := uc.repos.Workspace().GetDefault(txctx, req.Account.Sub)
		if err == nil {
			return &ProjectSetupDefaultRes{WorkspaceId: existing.Id, WorkspaceTier: existing.Tier.Name}, nil
		}

		if !errors.Is(err, database.ErrRecordNotFound) {
			return nil, err
		}

		// if there is a record not found error, we should create a new one
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

		ws, err := uc.repos.Workspace().Create(txctx, entity)
		if err != nil {
			return nil, err
		}

		res := &ProjectSetupDefaultRes{WorkspaceId: ws.Id, WorkspaceTier: ws.Tier.Name}
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	// must clear the cache because of new workspace
	cacheKey := cache.Key("WORKSPACES_OF_ACCOUNT", req.Account.Sub)
	if err := uc.cache.Del(cacheKey); err != nil {
		return nil, err
	}
	return res.(*ProjectSetupDefaultRes), err
}
