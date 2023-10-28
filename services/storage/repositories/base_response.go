package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Response interface {
	Create(ctx context.Context, docs []entities.Response) ([]entities.TSEntity, error)
}
