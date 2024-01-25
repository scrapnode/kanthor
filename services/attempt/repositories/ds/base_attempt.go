package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Attempt interface {
	Scan(ctx context.Context, query *entities.ScanningQuery, next int64, count int) chan *entities.ScanningResult[[]entities.Attempt]
	ListRequests(ctx context.Context, attempts map[string]*entities.Attempt) (map[string]*entities.Request, error)
	Update(ctx context.Context, updates map[string]*entities.AttemptState) map[string]error
}
