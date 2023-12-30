package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type MessageResponsetMaps struct {
	Maps    map[string][]entities.Response
	Success map[string]string
}

type Response interface {
	ListMessages(ctx context.Context, epId string, msgIds []string) (*MessageResponsetMaps, error)
	List(ctx context.Context, epId string, query *entities.ScanningQuery) ([]entities.Response, error)
	Get(ctx context.Context, epId, msgId, id string) (*entities.Response, error)
}
