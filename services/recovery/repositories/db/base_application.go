package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Application interface {
	Scan(ctx context.Context, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Application]
}
