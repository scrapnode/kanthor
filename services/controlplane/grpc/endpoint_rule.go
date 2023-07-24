package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/transform"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

type endpointRule struct {
	protos.UnimplementedEndpointRuleServer
	service *controlplane
	pipe    pipeline.Middleware
}

func (server *endpointRule) List(ctx context.Context, req *protos.EndpointRuleListReq) (*protos.EndpointRuleListRes, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.EndpointRule().List(ctx, request.(*usecase.EndpointRuleListReq))
		return
	})

	response, err := run(ctx, transform.EndpointRuleListReq(ctx, req))
	if err != nil {
		return nil, err
	}

	return transform.EndpointRuleListRes(ctx, response.(*usecase.EndpointRuleListRes)), nil
}

func (server *endpointRule) Get(ctx context.Context, req *protos.EndpointRuleGetReq) (*protos.EndpointRuleEntity, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.EndpointRule().Get(ctx, request.(*usecase.EndpointRuleGetReq))
		return
	})

	response, err := run(ctx, transform.EndpointRuleGetReq(ctx, req))
	if err != nil {
		return nil, err
	}

	return transform.EndpointRuleGetRes(ctx, response.(*usecase.EndpointRuleGetRes)), nil
}
