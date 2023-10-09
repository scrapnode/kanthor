package repos

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Msg struct {
	AppId     string `json:"app_id"`
	Id        string `json:"id"`
	Tier      string `json:"tier"`
	Timestamp int64  `json:"timestamp"`
}

type Message interface {
	Scan(ctx context.Context, appId string, from, to time.Time) (map[string]Msg, error)
	ListByIds(ctx context.Context, ids []string) ([]entities.Message, error)
}
