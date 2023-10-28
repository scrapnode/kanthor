package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Application interface {
	BulkCreate(ctx context.Context, docs []entities.Application) ([]string, error)
}
