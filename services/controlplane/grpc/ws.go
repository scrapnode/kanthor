package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ws struct {
	protos.UnimplementedWsServer
	service *controlplane
}

func (server *ws) Get(ctx context.Context, req *protos.WsGetReq) (*protos.Workspace, error) {
	request := &usecase.WorkspaceGetReq{Id: req.Id}

	response, err := server.service.uc.Workspace().Get(ctx, request)
	if err != nil {
		server.service.logger.Errorw(err.Error(), "workspace_id", req.Id)
		return nil, status.Error(codes.Internal, "unable to get workspace")
	}

	res := &protos.Workspace{
		Id:        response.Workspace.Id,
		CreatedAt: response.Workspace.CreatedAt,
		UpdatedAt: response.Workspace.UpdatedAt,
		DeletedAt: response.Workspace.DeletedAt,
		OwnerId:   response.Workspace.OwnerId,
		Name:      response.Workspace.Name,
		Tier: &protos.WorkspaceTier{
			WorkspaceId: response.Workspace.Tier.WorkspaceId,
			Name:        response.Workspace.Tier.Name,
		},
	}
	return res, nil
}
