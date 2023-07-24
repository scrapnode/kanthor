package transform

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

func Account(acc *authenticator.Account) *protos.AccountEntity {
	return &protos.AccountEntity{
		Sub:         acc.Sub,
		Name:        acc.Name,
		Picture:     acc.Picture,
		Email:       acc.Email,
		PhoneNumber: acc.PhoneNumber,
	}
}

func WorkspaceListOfAccountRes(ctx context.Context, res *usecase.WorkspaceListOfAccountRes) *protos.AccountListWorkspacesRes {
	returning := &protos.AccountListWorkspacesRes{Data: []*protos.WorkspaceEntity{}}
	for _, ws := range res.Workspaces {
		returning.Data = append(returning.Data, Workspace(&ws))
	}
	return returning
}
