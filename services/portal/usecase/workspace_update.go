package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type WorkspaceUpdateReq struct {
	Id   string
	Name string
}

func (req *WorkspaceUpdateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("id", req.Id, entities.IdNsWs),
		validator.StringRequired("name", req.Name),
	)
}

type WorkspaceUpdateRes struct {
	Doc *entities.Workspace
}

func (uc *workspace) Update(ctx context.Context, req *WorkspaceUpdateReq) (*WorkspaceUpdateRes, error) {
	ws, err := uc.repositories.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		ws, err := uc.repositories.Workspace().Get(txctx, req.Id)
		if err != nil {
			return nil, err
		}

		ws.Name = req.Name
		ws.SetAT(uc.infra.Timer.Now())
		return uc.repositories.Workspace().Update(txctx, ws)
	})
	if err != nil {
		return nil, err
	}

	res := &WorkspaceUpdateRes{Doc: ws.(*entities.Workspace)}
	return res, nil
}
