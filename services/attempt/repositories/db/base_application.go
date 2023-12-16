package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Application interface {
	Get(ctx context.Context, id string) (*entities.ApplicationWithRelationship, error)
	Scan(ctx context.Context, size int, cursor string) ([]entities.Application, error)
	GetTiers(ctx context.Context, apps []entities.Application) (map[string]string, error)
}
