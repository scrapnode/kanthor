package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
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

	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	request := &usecase.EndpointListReq{
		Workspace: ws,
		AppId:     req.AppId,
		ListReq:   structure.ListReq{Cursor: req.Cursor, Search: req.Search, Limit: int(req.Limit), Ids: req.Ids},
	}
	response, err := run(ctx, request)
	if err != nil {
		return nil, err
	}

	// transformation
	cast := response.(*usecase.EndpointListRes)
	res := &protos.EndpointListRes{Cursor: cast.Cursor, Data: []*protos.EndpointEntity{}}
	for _, ep := range cast.Data {
		endpoint := &protos.EndpointEntity{
			Id:        ep.Id,
			CreatedAt: ep.CreatedAt,
			UpdatedAt: ep.UpdatedAt,
			AppId:     ep.AppId,
			Name:      ep.Name,
			Method:    ep.Method,
			Uri:       ep.Uri,
		}
		res.Data = append(res.Data, endpoint)
	}

	return res, nil
}

func (server *endpoint) Get(ctx context.Context, req *protos.EndpointGetReq) (*protos.EndpointEntity, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Endpoint().Get(ctx, request.(*usecase.EndpointGetReq))
		return
	})

	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	request := &usecase.EndpointGetReq{Workspace: ws, AppId: req.AppId, Id: req.Id}
	response, err := run(ctx, request)
	if err != nil {
		return nil, err
	}

	// transformation
	cast := response.(*usecase.EndpointGetRes)
	res := &protos.EndpointEntity{
		Id:        cast.Endpoint.Id,
		CreatedAt: cast.Endpoint.CreatedAt,
		UpdatedAt: cast.Endpoint.UpdatedAt,
		AppId:     cast.Endpoint.AppId,
		Name:      cast.Endpoint.Name,
		Method:    cast.Endpoint.Method,
		Uri:       cast.Endpoint.Uri,
	}
	return res, nil
}
