package usecase

import (
	"context"
	"net/http"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type EndpointCreateIn struct {
	WsId  string
	AppId string
	Name  string

	SecretKey string
	Uri       string
	Method    string
}

func (in *EndpointCreateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.StringRequired("name", in.Name),
		validator.StringRequired("secret_key", in.SecretKey),
		validator.StringLen("secret_key", in.SecretKey, 16, 32),
		validator.StringUri("uri", in.Uri),
		validator.StringOneOf("method", in.Method, []string{http.MethodPost, http.MethodPut}),
	)
}

type EndpointCreateOut struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Create(ctx context.Context, in *EndpointCreateIn) (*EndpointCreateOut, error) {
	app, err := uc.repositories.Application().Get(ctx, in.WsId, in.AppId)
	if err != nil {
		return nil, err
	}

	doc := &entities.Endpoint{
		AppId:     app.Id,
		Name:      in.Name,
		SecretKey: in.SecretKey,
		Method:    in.Method,
		Uri:       in.Uri,
	}
	doc.GenId()
	doc.SetAT(uc.infra.Timer.Now())
	doc.GenSecretKey()

	ep, err := uc.repositories.Endpoint().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	return &EndpointCreateOut{Doc: ep}, nil
}
