package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Application interface {
	CreateBatch(ctx context.Context, docs []entities.Application) ([]string, error)
	Count(ctx context.Context, wsId string, query *entities.PagingQuery) (int64, error)
	Get(ctx context.Context, wsId, id string) (*entities.Application, error)
}
