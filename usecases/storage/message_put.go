package storage

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type MessagePutReq struct {
	Docs []entities.Message
}

type MessagePutRes struct {
	Entities []entities.Entity
}

func (uc *message) Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error) {
	entities, err := uc.repos.Message().Create(ctx, req.Docs)
	if err != nil {
		return nil, err
	}

	res := &MessagePutRes{Entities: entities}
	return res, nil
}
