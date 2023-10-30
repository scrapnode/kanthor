package repositories

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Attempt interface {
	Create(ctx context.Context, docs []entities.Attempt) ([]string, error)
	Scan(ctx context.Context, from, to time.Time, matching int64) ([]entities.Attempt, error)
}
