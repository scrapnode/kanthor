package transform

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

func Endpoint(ep *entities.Endpoint) *protos.EndpointEntity {
	return &protos.EndpointEntity{
		Id:        ep.Id,
		CreatedAt: ep.CreatedAt,
		UpdatedAt: ep.UpdatedAt,
		AppId:     ep.AppId,
		Name:      ep.Name,
		Method:    ep.Method,
		Uri:       ep.Uri,
	}
}

func EndpointListReq(ctx context.Context, req *protos.EndpointListReq) *usecase.EndpointListReq {
	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	return &usecase.EndpointListReq{
		Workspace: ws,
		AppId:     req.AppId,
		ListReq: structure.ListReq{
			Cursor: req.Cursor,
			Search: req.Search,
			Limit:  int(req.Limit),
			Ids:    req.Ids,
		},
	}
}

func EndpointListRes(ctx context.Context, res *usecase.EndpointListRes) *protos.EndpointListRes {
	returning := &protos.EndpointListRes{Cursor: res.Cursor, Data: []*protos.EndpointEntity{}}
	for _, ep := range res.Data {
		returning.Data = append(returning.Data, Endpoint(&ep))
	}
	return returning
}

func EndpointGetReq(ctx context.Context, req *protos.EndpointGetReq) *usecase.EndpointGetReq {
	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	return &usecase.EndpointGetReq{
		Workspace: ws,
		AppId:     req.AppId,
		Id:        req.Id,
	}
}

func EndpointGetRes(ctx context.Context, res *usecase.EndpointGetRes) *protos.EndpointEntity {
	return Endpoint(res.Endpoint)
}
