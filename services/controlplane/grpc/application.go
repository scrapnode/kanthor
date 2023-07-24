package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

type application struct {
	protos.UnimplementedApplicationServer
	service *controlplane
	pipe    pipeline.Middleware
}

func (server *application) List(ctx context.Context, req *protos.ApllicationListReq) (*protos.ApllicationListRes, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Application().List(ctx, request.(*usecase.ApplicationListReq))
		return
	})

	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	request := &usecase.ApplicationListReq{
		Workspace: ws,
		ListReq:   structure.ListReq{Cursor: req.Cursor, Search: req.Search, Limit: int(req.Limit), Ids: req.Ids},
	}
	response, err := run(ctx, request)
	if err != nil {
		return nil, err
	}

	// transformation
	cast := response.(*usecase.ApplicationListRes)
	res := &protos.ApllicationListRes{Cursor: cast.Cursor, Data: []*protos.ApplicationEntity{}}
	for _, app := range cast.Data {
		res.Data = append(res.Data, &protos.ApplicationEntity{
			Id:          app.Id,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
			WorkspaceId: app.WorkspaceId,
			Name:        app.Name,
		})
	}

	return res, nil
}

func (server *application) Get(ctx context.Context, req *protos.ApllicationGetReq) (*protos.ApplicationEntity, error) {
	run := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Application().Get(ctx, request.(*usecase.ApplicationGetReq))
		return
	})

	ws := ctx.Value(usecase.CtxWorkspace).(*entities.Workspace)
	request := &usecase.ApplicationGetReq{Workspace: ws, Id: req.Id}
	response, err := run(ctx, request)
	if err != nil {
		return nil, err
	}

	// transformation
	app := response.(*usecase.ApplicationGetRes).Application
	res := &protos.ApplicationEntity{
		Id:          app.Id,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
		WorkspaceId: app.WorkspaceId,
		Name:        app.Name,
	}
	return res, nil
}
