package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type workspace struct {
	protos.UnimplementedWorkspaceServer
	service *controlplane
}

func (server *workspace) Get(ctx context.Context, req *protos.WorkspaceGetReq) (*protos.IWorkspace, error) {
	request := &usecase.WorkspaceGetReq{Id: req.Id}
	response, err := server.service.uc.Workspace().Get(ctx, request)
	if err != nil {
		server.service.logger.Error(err.Error(), "request", utils.Stringify(req))
		return nil, status.Error(codes.Internal, "oops, something went wrong")
	}

	res := &protos.IWorkspace{
		Id:        response.Workspace.Id,
		CreatedAt: response.Workspace.CreatedAt,
		UpdatedAt: response.Workspace.UpdatedAt,
		OwnerId:   response.Workspace.OwnerId,
		Name:      response.Workspace.Name,
		Tier: &protos.IWorkspaceTier{
			WorkspaceId: response.Workspace.Tier.WorkspaceId,
			Name:        response.Workspace.Tier.Name,
		},
	}
	return res, nil
}
