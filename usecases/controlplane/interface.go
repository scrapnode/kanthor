package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Controlplane interface {
	patterns.Connectable
	Workspace() Workspace
}

type Workspace interface {
	Get(ctx context.Context, req *WorkspaceGetReq) (*WorkspaceGetRes, error)
	ListOfAccount(ctx context.Context, req *WorkspaceListOfAccountReq) (*WorkspaceListOfAccountRes, error)
	GetByAccount(ctx context.Context, req *WorkspaceGetByAccountReq) (*WorkspaceGetByAccountRes, error)
}

type WorkspaceGetReq struct {
	Id string
}

type WorkspaceGetRes struct {
	Workspace *entities.Workspace
}

type WorkspaceListOfAccountReq struct {
	structure.ListReq
	AccountSub string
}

type WorkspaceListOfAccountRes struct {
	Workspaces []entities.Workspace
}

type WorkspaceGetByAccountReq struct {
	WorkspaceId string
	AccountSub  string
}

type WorkspaceGetByAccountRes struct {
	Workspace *entities.Workspace
	Privilege *entities.WorkspacePrivilege
}
