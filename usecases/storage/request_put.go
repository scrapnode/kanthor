package storage

import (
	"context"
)

func (uc *request) Put(ctx context.Context, req *RequestPutReq) (*RequestPutRes, error) {
	entities, err := uc.repos.Request().Create(ctx, req.Docs)
	if err != nil {
		return nil, err
	}

	res := &RequestPutRes{Entities: entities}
	return res, nil
}
