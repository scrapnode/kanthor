package controlplane

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Controlplane interface {
	patterns.Connectable
	Workspace() Workspace
}

type Workspace interface {
	ListOfAccount(ctx context.Context, req *WorkspaceListOfAccountReq) (*WorkspaceListOfAccountRes, error)
	GetByAccount(ctx context.Context, req *WorkspaceGetByAccountReq) (*WorkspaceGetByAccountRes, error)
}

type WorkspaceGetReq struct {
	Id string
}

type WorkspaceGetRes struct {
	Workspace *entities.Workspace
}

type WorkspaceGetByAccountReq struct {
	WorkspaceId string
	Account     *authenticator.Account
}

type WorkspaceGetByAccountRes struct {
	Workspace *entities.Workspace
}

type WorkspaceListOfAccountReq struct {
	structure.ListReq
	Account *authenticator.Account
}

type WorkspaceListOfAccountRes struct {
	Workspaces []entities.Workspace
}
