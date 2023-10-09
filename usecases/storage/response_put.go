package storage

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type ResponsePutReq struct {
	Docs []entities.Response
}

type ResponsePutRes struct {
	Entities []entities.Entity
}

func (uc *response) Put(ctx context.Context, req *ResponsePutReq) (*ResponsePutRes, error) {
	entities, err := uc.repos.Response().Create(ctx, req.Docs)
	if err != nil {
		return nil, err
	}

	res := &ResponsePutRes{Entities: entities}
	return res, nil
}
