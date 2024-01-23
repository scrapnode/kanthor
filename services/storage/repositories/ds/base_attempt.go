package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Attempt interface {
	Create(ctx context.Context, docs []*entities.Attempt) ([]string, error)
}
