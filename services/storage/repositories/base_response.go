package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
)

type Response interface {
	Create(ctx context.Context, docs []*entities.Response) ([]string, error)
}
