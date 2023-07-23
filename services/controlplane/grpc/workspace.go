package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

type workspace struct {
	protos.UnimplementedWorkspaceServer
	service *controlplane
	pipe    pipeline.Middleware
}

func (server *workspace) Get(ctx context.Context, req *protos.WorkspaceGetReq) (*protos.IWorkspace, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Workspace().Get(ctx, request.(*usecase.WorkspaceGetReq))
		return
	})

	request := &usecase.WorkspaceGetReq{Id: req.Id}
	response, err := run(ctx, request)
	if err != nil {
		return nil, err
	}

	// transformation
	ws := response.(*usecase.WorkspaceGetRes).Workspace
	res := &protos.IWorkspace{
		Id:        ws.Id,
		CreatedAt: ws.CreatedAt,
		UpdatedAt: ws.UpdatedAt,
		OwnerId:   ws.OwnerId,
		Name:      ws.Name,
		Tier: &protos.IWorkspaceTier{
			WorkspaceId: ws.Tier.WorkspaceId,
			Name:        ws.Tier.Name,
		},
	}
	return res, nil
}
