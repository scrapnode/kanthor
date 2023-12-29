package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Request interface {
	List(ctx context.Context, epId string, query *entities.ScanningQuery) ([]entities.Request, error)
	Get(ctx context.Context, epId, msgId, id string) (*entities.Request, error)
}
