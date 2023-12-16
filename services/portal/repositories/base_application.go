package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Application interface {
	BulkCreate(ctx context.Context, docs []entities.Application) ([]string, error)
}
