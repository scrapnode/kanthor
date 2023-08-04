package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

func (uc *endpoint) Create(ctx context.Context, req *EndpointCreateReq) (*EndpointCreateRes, error) {
	doc := &entities.Endpoint{
		AppId:     req.AppId,
		Name:      req.Name,
		SecretKey: req.SecretKey,
		Method:    req.Method,
		Uri:       req.Uri,
	}
	doc.GenId()
	doc.SetAT(uc.timer.Now())
	doc.GenSecretKey()

	ep, err := uc.repos.Endpoint().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	res := &EndpointCreateRes{Doc: ep}
	return res, nil
}
