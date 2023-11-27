package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
)

type Request interface {
	Create(ctx context.Context, docs []*entities.Request) ([]string, error)
}
