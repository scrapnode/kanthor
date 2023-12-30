package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type MessageRequestMaps struct {
	Maps   map[string][]entities.Request
	MsgIds []string
}

type Request interface {
	ListMessages(ctx context.Context, epId string, query *entities.ScanningQuery) (*MessageRequestMaps, error)
	GetMessage(ctx context.Context, epId, msgId string) (*MessageRequestMaps, error)
	List(ctx context.Context, epId string, query *entities.ScanningQuery) ([]entities.Request, error)
}
