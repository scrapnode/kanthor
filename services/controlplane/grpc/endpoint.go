package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/transform"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

type endpoint struct {
	protos.UnimplementedEndpointServer
	service *controlplane
	pipe    pipeline.Middleware
}

func (server *endpoint) List(ctx context.Context, req *protos.EndpointListReq) (*protos.EndpointListRes, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Endpoint().List(ctx, request.(*usecase.EndpointListReq))
		return
	})

	response, err := run(ctx, transform.EndpointListReq(ctx, req))
	if err != nil {
		return nil, err
	}

	return transform.EndpointListRes(ctx, response.(*usecase.EndpointListRes)), nil
}

func (server *endpoint) Get(ctx context.Context, req *protos.EndpointGetReq) (*protos.EndpointEntity, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Endpoint().Get(ctx, request.(*usecase.EndpointGetReq))
		return
	})

	response, err := run(ctx, transform.EndpointGetReq(ctx, req))
	if err != nil {
		return nil, err
	}

	return transform.EndpointGetRes(ctx, response.(*usecase.EndpointGetRes)), nil
}
