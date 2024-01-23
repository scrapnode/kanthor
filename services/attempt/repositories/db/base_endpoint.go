package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Endpoint interface {
	Scan(ctx context.Context, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Endpoint]
}
