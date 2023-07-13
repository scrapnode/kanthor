package grpc

import (
	"context"
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

func (server *account) ListWorkspaces(ctx context.Context, _ *protos.ListWorkspacesReq) (*protos.ListWorkspacesRes, error) {
	wsIds := authenticator.WorkspaceIdsFromContext(ctx)
	request := &usecase.WorkspaceListByIdsReq{Ids: wsIds}
	response, err := server.service.uc.Workspace().ListByIds(ctx, request)
	if err != nil {
		server.service.logger.Errorw(err.Error(), "ws_ids", wsIds)
		return nil, status.Error(codes.Internal, "unable to retrieve data from list workspace")
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
