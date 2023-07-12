package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Workspace interface {
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	ListByIds(ctx context.Context, ids []string) ([]entities.Workspace, error)
}
