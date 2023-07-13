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
	ListByIds(ctx context.Context, req *WorkspaceListByIdsReq) (*WorkspaceListByIdsRes, error)
}

type WorkspaceGetReq struct {
	Id string
}

type WorkspaceGetRes struct {
	Workspace entities.Workspace
}

type WorkspaceListByIdsReq struct {
	Ids []string
}

type WorkspaceListByIdsRes struct {
	Workspaces []entities.Workspace
}
