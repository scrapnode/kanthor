package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
)

type Workspace interface {
	Get(ctx context.Context, id string) (*entities.Workspace, error)
	ListOfAccountSub(ctx context.Context, sub string, opts ...structure.ListOps) (*structure.ListRes[entities.Workspace], error)
	GetByAccountSub(ctx context.Context, id, sub string) (*entities.Workspace, error)
}
