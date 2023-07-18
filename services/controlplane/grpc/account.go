package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type account struct {
	protos.UnimplementedAccountServer
	service *controlplane
}

func (server *account) Get(ctx context.Context, req *protos.AccountGetReq) (*protos.IAccount, error) {
	acc := ctx.Value(authenticator.CtxAuthAccount).(*authenticator.Account)

	res := &protos.IAccount{
		Sub:         acc.Sub,
		Name:        acc.Name,
		Picture:     acc.Picture,
		Email:       acc.Email,
		PhoneNumber: acc.PhoneNumber,
	}
	return res, nil
}

func (server *account) ListWorkspaces(ctx context.Context, req *protos.AccountListWorkspacesReq) (*protos.AccountListWorkspacesRes, error) {
	acc := ctx.Value(authenticator.CtxAuthAccount).(*authenticator.Account)

	request := &usecase.WorkspaceListOfAccountReq{Account: acc}
	response, err := server.service.uc.Workspace().ListOfAccount(ctx, request)
	if err != nil {
		server.service.logger.Errorw(err.Error(), "request", utils.Stringify(req))
		return nil, status.Error(codes.Internal, "oops, something went wrong")
	}

	res := &protos.AccountListWorkspacesRes{Data: []*protos.IWorkspace{}}
	for _, workspace := range response.Workspaces {
		res.Data = append(res.Data, &protos.IWorkspace{
			Id:        workspace.Id,
			CreatedAt: workspace.CreatedAt,
			UpdatedAt: workspace.UpdatedAt,
			OwnerId:   workspace.OwnerId,
			Name:      workspace.Name,
			Tier: &protos.IWorkspaceTier{
				WorkspaceId: workspace.Tier.WorkspaceId,
				Name:        workspace.Tier.Name,
			},
		})
	}

	return res, nil
}
