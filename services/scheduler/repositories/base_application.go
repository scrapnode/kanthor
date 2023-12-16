package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Application interface {
	Get(ctx context.Context, id string) (*entities.Application, error)
}
