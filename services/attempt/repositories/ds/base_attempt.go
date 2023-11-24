package ds

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Attempt interface {
	Create(ctx context.Context, docs []entities.Attempt) ([]string, error)

	Count(ctx context.Context, appId string, from, to time.Time, next int64) (int64, error)
	Scan(ctx context.Context, from, to time.Time, next int64, limit int) chan *ScanResults[[]entities.Attempt]
	MarkComplete(ctx context.Context, reqId string, res *entities.Response) error
	MarkReschedule(ctx context.Context, reqId string, ts int64) error
	MarkIgnore(ctx context.Context, reqIds []string) error
}
