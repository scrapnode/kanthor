package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationUpdateReq struct {
	WorkspaceId string
	Id          string
	Name        string
}

func (req *ApplicationUpdateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("workspace_id", req.WorkspaceId, entities.IdNsWs),
		validator.StringStartsWith("id", req.Id, entities.IdNsApp),
		validator.StringRequired("name", req.Name),
	)
}

type ApplicationUpdateRes struct {
	Doc *entities.Application
}

func (uc *application) Update(ctx context.Context, req *ApplicationUpdateReq) (*ApplicationUpdateRes, error) {
	app, err := uc.repos.Transaction(ctx, func(txctx context.Context) (interface{}, error) {
		app, err := uc.repos.Application().Get(txctx, req.WorkspaceId, req.Id)
		if err != nil {
			return nil, err
		}

		app.Name = req.Name
		app.SetAT(uc.infra.Timer.Now())
		return uc.repos.Application().Update(txctx, app)
	})
	if err != nil {
		return nil, err
	}

	res := &ApplicationUpdateRes{Doc: app.(*entities.Application)}
	return res, nil
}
