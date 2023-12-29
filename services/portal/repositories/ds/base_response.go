package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Response interface {
	List(ctx context.Context, epId string, query *entities.ScanningQuery) ([]entities.Response, error)
	Get(ctx context.Context, epId, msgId string) ([]entities.Response, error)
}
