package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/transform"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

type application struct {
	protos.UnimplementedApplicationServer
	service *controlplane
	pipe    pipeline.Middleware
}

func (server *application) List(ctx context.Context, req *protos.ApplicationListReq) (*protos.ApplicationListRes, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Application().List(ctx, request.(*usecase.ApplicationListReq))
		return
	})

	response, err := run(ctx, transform.ApplicationListReq(ctx, req))
	if err != nil {
		return nil, err
	}

	return transform.ApplicationListRes(ctx, response.(*usecase.ApplicationListRes)), nil
}

func (server *application) Get(ctx context.Context, req *protos.ApplicationGetReq) (*protos.ApplicationEntity, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Application().Get(ctx, request.(*usecase.ApplicationGetReq))
		return
	})

	response, err := run(ctx, transform.ApplicationGetReq(ctx, req))
	if err != nil {
		return nil, err
	}

	return transform.ApplicationGetRes(ctx, response.(*usecase.ApplicationGetRes)), nil
}
