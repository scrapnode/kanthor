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

func (server *workspace) Get(ctx context.Context, req *protos.WorkspaceGetReq) (*protos.WorkspaceEntity, error) {
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
	cast := response.(*usecase.WorkspaceGetRes)
	res := &protos.WorkspaceEntity{
		Id:        cast.Workspace.Id,
		CreatedAt: cast.Workspace.CreatedAt,
		UpdatedAt: cast.Workspace.UpdatedAt,
		OwnerId:   cast.Workspace.OwnerId,
		Name:      cast.Workspace.Name,
		Tier: &protos.WorkspaceTierEntity{
			WorkspaceId: cast.Workspace.Tier.WorkspaceId,
			Name:        cast.Workspace.Tier.Name,
		},
	}
	return res, nil
}
