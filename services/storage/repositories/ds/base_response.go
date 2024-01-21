package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Response interface {
	Create(ctx context.Context, docs []*entities.Response) ([]string, error)
}
