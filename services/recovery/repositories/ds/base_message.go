package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Message interface {
	Scan(ctx context.Context, appId string, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Message]
}
