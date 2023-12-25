package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Message interface {
	List(ctx context.Context, appId string, query *entities.ScanningQuery) ([]entities.Message, error)
	Get(ctx context.Context, appId, id string) (*entities.Message, error)
}
