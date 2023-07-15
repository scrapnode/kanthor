package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Workspace interface {
	List(ctx context.Context, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error)
	Get(ctx context.Context, id string) (*entities.Workspace, error)
}