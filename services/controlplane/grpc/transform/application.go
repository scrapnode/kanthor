package transform

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

func Application(app *entities.Application) *protos.ApplicationEntity {
	return &protos.ApplicationEntity{
		Id:          app.Id,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
		WorkspaceId: app.WorkspaceId,
		Name:        app.Name,
	}
}

func ApplicationListReq(ctx context.Context, req *protos.ApplicationListReq) *usecase.ApplicationListReq {
	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)

	return &usecase.ApplicationListReq{
		Workspace: ws,
		ListReq: structure.ListReq{
			Cursor: req.Cursor,
			Search: req.Search,
			Limit:  int(req.Limit),
			Ids:    req.Ids,
		},
	}
}

func ApplicationListRes(ctx context.Context, res *usecase.ApplicationListRes) *protos.ApplicationListRes {
	returning := &protos.ApplicationListRes{Cursor: res.Cursor, Data: []*protos.ApplicationEntity{}}
	for _, app := range res.Data {
		returning.Data = append(returning.Data, Application(&app))
	}
	return returning
}

func ApplicationGetReq(ctx context.Context, req *protos.ApplicationGetReq) *usecase.ApplicationGetReq {
	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	return &usecase.ApplicationGetReq{
		Workspace: ws,
		Id:        req.Id,
	}
}

func ApplicationGetRes(ctx context.Context, res *usecase.ApplicationGetRes) *protos.ApplicationEntity {
	return Application(res.Application)
}
