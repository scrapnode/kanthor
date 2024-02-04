package usecase

import (
	"context"
	"errors"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/telemetry"
	"go.opentelemetry.io/otel/trace"
)

type WorkspaceGetIn struct {
	AccId string
	Id    string
}

func (in *WorkspaceGetIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("acc_id", in.AccId),
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
	)
}

type WorkspaceGetOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error) {
	subctx, span := ctx.Value(telemetry.CtxTracer).(trace.Tracer).Start(ctx, "usecase.workspace.get")
	defer func() {
		span.End()
	}()

	ws, err := uc.repositories.Database().Workspace().Get(subctx, in.Id)
	if err != nil {
		return nil, err
	}

	isOwner := ws.OwnerId == in.AccId
	if isOwner {
		return &WorkspaceGetOut{ws}, nil
	}

	wsc, err := uc.repositories.Database().WorkspaceCredentials().Get(subctx, in.AccId)
	if err != nil {
		return nil, err
	}

	isWorker := wsc.WsId == in.Id
	if isWorker {
		return &WorkspaceGetOut{ws}, nil
	}

	return nil, errors.New("SDK.USECASE.WORKSPACE.GET.NOT_OWN.ERROR")
}
