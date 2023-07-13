package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type account struct {
	protos.UnimplementedAccountServer
	service *controlplane
}

func (server *account) ListWorkspaces(ctx context.Context, req *protos.ListWorkspacesReq) (*protos.ListWorkspacesRes, error) {
	acc := authenticator.AccountFromContext(ctx)
	request := &usecase.WorkspaceListOfAccountReq{
		ListReq: structure.ListReq{
			Cursor: req.Cursor,
			Search: req.Search,
			Limit:  int(req.Limit),
		},
		AccountSub: acc.Sub,
	}

	response, err := server.service.uc.Workspace().ListOfAccount(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &protos.ListWorkspacesRes{Data: []*protos.Workspace{}}
	for _, workspace := range response.Workspaces {
		res.Data = append(res.Data, &protos.Workspace{
			Id:        workspace.Id,
			CreatedAt: workspace.CreatedAt,
			UpdatedAt: workspace.UpdatedAt,
			DeletedAt: workspace.DeletedAt,
			OwnerId:   workspace.OwnerId,
			Name:      workspace.Name,
			Tier: &protos.WorkspaceTier{
				WorkspaceId: workspace.Tier.WorkspaceId,
				Name:        workspace.Tier.Name,
			},
		})
	}

	return res, nil
}
