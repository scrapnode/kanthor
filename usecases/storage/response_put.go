package storage

import (
	"context"
)

func (uc *response) Put(ctx context.Context, req *ResponsePutReq) (*ResponsePutRes, error) {
	entities, err := uc.repos.Response().Create(ctx, req.Docs)
	if err != nil {
		return nil, err
	}

	res := &ResponsePutRes{Entities: entities}
	return res, nil
}
