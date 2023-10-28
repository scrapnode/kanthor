package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceGetReq struct {
	Id string
}

func (req *WorkspaceGetReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("id", req.Id, entities.IdNsWs),
	)
}

type WorkspaceGetRes struct {
	Workspace *entities.Workspace
}

func (uc *workspace) Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error) {
	key := utils.Key("portal", req.Id)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*WorkspaceGetRes, error) {
		ws, err := uc.repos.Workspace().Get(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		res := &WorkspaceGetRes{Workspace: ws}
		return res, nil
	})
}
