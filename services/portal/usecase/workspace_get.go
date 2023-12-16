package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceGetIn struct {
	Id string
}

func (in *WorkspaceGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
	)
}

type WorkspaceGetOut struct {
	Workspace *entities.Workspace
}

func (uc *workspace) Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error) {
	key := utils.Key("portal", in.Id)
	return cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*WorkspaceGetOut, error) {
		ws, err := uc.repositories.Workspace().Get(ctx, in.Id)
		if err != nil {
			return nil, err
		}

		return &WorkspaceGetOut{Workspace: ws}, nil
	})
}
