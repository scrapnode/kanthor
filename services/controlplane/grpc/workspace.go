package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/transform"
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

	response, err := run(ctx, transform.WorkspaceGetReq(ctx, req))
	if err != nil {
		return nil, err
	}

	return transform.WorkspaceGetRes(ctx, response.(*usecase.WorkspaceGetRes)), nil
}
