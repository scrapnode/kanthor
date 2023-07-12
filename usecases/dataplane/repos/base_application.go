package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Application interface {
	GetWithWorkspace(ctx context.Context, id string) (*ApplicationWithWorkspace, error)
}

type ApplicationWithWorkspace struct {
	entities.Application
	Workspace entities.Workspace
}
