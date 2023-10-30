package repositories

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Attempt interface {
	BulkCreate(ctx context.Context, docs []entities.Attempt) ([]string, error)
	Scan(ctx context.Context, from, to time.Time, matching int64) ([]entities.Attempt, error)
	MarkComplete(ctx context.Context, reqId string, res *entities.Response) error
	MarkReschedule(ctx context.Context, reqId string, ts int64) error
}
