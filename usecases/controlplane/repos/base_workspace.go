package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Workspace interface {
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	ListByIds(ctx context.Context, ids []string) (*structure.ListRes[entities.Workspace], error)
}
