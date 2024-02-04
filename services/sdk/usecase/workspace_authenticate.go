package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/telemetry"
	"go.opentelemetry.io/otel/trace"
)

type WorkspaceAuthenticateIn struct {
	User string
	Pass string
}

func (in *WorkspaceAuthenticateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("user", in.User, entities.IdNsWsc),
		validator.StringRequired("pass", in.Pass),
	)
}

type WorkspaceAuthenticateOut struct {
	Credentials *entities.WorkspaceCredentials
}

func (uc *workspace) Authenticate(ctx context.Context, in *WorkspaceAuthenticateIn) (*WorkspaceAuthenticateOut, error) {
	tracer := ctx.Value(telemetry.CtxTracer).(trace.Tracer)
	subctx, span := tracer.Start(ctx, "usecase.workspace.authenticate")
	defer func() {
		span.End()
	}()

	credentials, err := uc.repositories.Database().WorkspaceCredentials().Get(subctx, in.User)
	if err != nil {
		return nil, err
	}

	expired := credentials.ExpiredAt > 0 && credentials.ExpiredAt < uc.infra.Timer.Now().UnixMilli()
	if expired {
		expiredAt := time.UnixMilli(credentials.ExpiredAt).Format(time.RFC3339)
		return nil, fmt.Errorf("workspace credentials was expired (%s)", expiredAt)
	}

	_, subspan := tracer.Start(ctx, "usecase.workspace.authenticate.password.compare")
	if err := utils.PasswordCompare(in.Pass, credentials.Hash); err != nil {
		subspan.End()
		return nil, err
	}
	subspan.End()

	return &WorkspaceAuthenticateOut{Credentials: credentials}, nil
}
