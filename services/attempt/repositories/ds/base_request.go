package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Request interface {
	Scan(ctx context.Context, epId string, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Request]
}
