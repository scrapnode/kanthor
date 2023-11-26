package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
)

type Application interface {
	Get(ctx context.Context, id string) (*entities.Application, error)
}
