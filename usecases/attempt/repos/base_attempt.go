package repos

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Attempt interface {
	Create(ctx context.Context, docs []entities.Attempt) ([]string, error)
}
