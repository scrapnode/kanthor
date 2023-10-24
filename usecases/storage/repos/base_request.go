package repos

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Request interface {
	Create(ctx context.Context, docs []entities.Request) ([]entities.TSEntity, error)
}
