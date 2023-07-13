package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Controlplane interface {
	patterns.Connectable
	Workspace() Workspace
}

type Workspace interface {
	Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error)
	ListByIds(ctx context.Context, req *WorkspaceListReq) (*WorkspaceListRes, error)
}

type WorkspaceGetReq struct {
	Id string
}

type WorkspaceGetRes struct {
	Workspace entities.Workspace
}

type WorkspaceListReq struct {
	Ids []string
}

type WorkspaceListRes struct {
	Workspaces []entities.Workspace
}
