package ds

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/internal/domain/entities"
)

type Message interface {
	Count(ctx context.Context, appId string, from, to time.Time) (int64, error)
	Scan(ctx context.Context, appId string, from, to time.Time, limit int) chan *ScanResults[map[string]entities.Message]
}
