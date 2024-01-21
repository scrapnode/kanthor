package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Request interface {
	Create(ctx context.Context, docs []*entities.Request) ([]string, error)
}
