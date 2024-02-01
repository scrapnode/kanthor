package usecase

import (
	"context"
	"errors"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceGetIn struct {
	AccId string
	Id    string
}

func (in *WorkspaceGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("acc_id", in.AccId),
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
	)
}

type WorkspaceGetOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error) {
	ws, err := uc.repositories.Database().Workspace().Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	isOwner := ws.OwnerId == in.AccId
	if isOwner {
		return &WorkspaceGetOut{ws}, nil
	}

	wsc, err := uc.repositories.Database().WorkspaceCredentials().Get(ctx, in.AccId)
	if err != nil {
		return nil, err
	}

	isWorker := wsc.WsId == in.Id
	if isWorker {
		return &WorkspaceGetOut{ws}, nil
	}

	return nil, errors.New("SDK.USECASE.WORKSPACE.GET.NOT_OWN.ERROR")
}
