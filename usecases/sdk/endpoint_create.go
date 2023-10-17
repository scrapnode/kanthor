package sdk

import (
	"context"
	"net/http"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointCreateReq struct {
	AppId string
	Name  string

	SecretKey string
	Uri       string
	Method    string
}

func (req *EndpointCreateReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
		validator.StringRequired("name", req.Name),
		validator.StringRequired("secret_key", req.SecretKey),
		validator.StringLen("secret_key", req.SecretKey, 16, 32),
		validator.StringUri("uri", req.Uri),
		validator.StringOneOf("method", req.Method, []string{http.MethodPost, http.MethodPut}),
	)
}

type EndpointCreateRes struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Create(ctx context.Context, req *EndpointCreateReq) (*EndpointCreateRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	app, err := uc.repos.Application().Get(ctx, ws, req.AppId)
	if err != nil {
		return nil, err
	}

	doc := &entities.Endpoint{
		AppId:     app.Id,
		Name:      req.Name,
		SecretKey: req.SecretKey,
		Method:    req.Method,
		Uri:       req.Uri,
	}
	doc.GenId()
	doc.SetAT(uc.infra.Timer.Now())
	doc.GenSecretKey()

	ep, err := uc.repos.Endpoint().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	res := &EndpointCreateRes{Doc: ep}
	return res, nil
}
