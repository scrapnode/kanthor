package sdk

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type ApplicationCreateReq struct {
	Name string
}

func (req *ApplicationCreateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("name", req.Name),
	)
}

type ApplicationCreateRes struct {
	Doc *entities.Application
}

func (uc *application) Create(ctx context.Context, req *ApplicationCreateReq) (*ApplicationCreateRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	doc := &entities.Application{WsId: ws.Id, Name: req.Name}
	doc.GenId()
	doc.SetAT(uc.infra.Timer.Now())

	app, err := uc.repos.Application().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	res := &ApplicationCreateRes{Doc: app}
	return res, nil
}
