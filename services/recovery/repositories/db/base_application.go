package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/routing"
)

type Application interface {
	Scan(ctx context.Context, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Application]
	GetRoutes(ctx context.Context, ids []string) (map[string][]routing.Route, error)
}
