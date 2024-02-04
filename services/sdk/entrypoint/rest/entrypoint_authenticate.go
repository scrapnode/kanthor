package rest

import (
	"context"

	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
	"github.com/scrapnode/kanthor/telemetry"
	"go.opentelemetry.io/otel/trace"
)

func RegisterWorkspaceResolver(uc usecase.Sdk) func(ctx context.Context, acc *authenticator.Account, id string) (*entities.Workspace, error) {
	return func(ctx context.Context, acc *authenticator.Account, id string) (*entities.Workspace, error) {
		subctx, span := ctx.Value(telemetry.CtxTracer).(trace.Tracer).Start(ctx, "entrypoint.authentication.workspace.resolve")
		defer func() {
			span.End()
		}()

		in := &usecase.WorkspaceGetIn{
			AccId: acc.Sub,
			Id:    id,
		}
		if err := in.Validate(); err != nil {
			return nil, err
		}

		out, err := uc.Workspace().Get(subctx, in)
		if err != nil {
			return nil, err
		}
		return out.Doc, nil
	}
}

var AuthzEngineInternal = "sdk.internal"

type internal struct {
	uc usecase.Sdk
}

func (verifier *internal) Verify(ctx context.Context, request *authenticator.Request) (*authenticator.Account, error) {
	user, pass, err := authenticator.ParseBasicCredentials(request.Credentials)
	if err != nil {
		return nil, err
	}

	in := &usecase.WorkspaceAuthenticateIn{User: user, Pass: pass}
	if err := in.Validate(); err != nil {
		return nil, err
	}

	out, err := verifier.uc.Workspace().Authenticate(ctx, in)
	if err != nil {
		return nil, err
	}

	account := &authenticator.Account{
		Sub:  out.Credentials.Id,
		Name: out.Credentials.Name,
		Metadata: map[string]string{
			gateway.MetaWorkspaceId: out.Credentials.WsId,
		},
	}
	return account, nil
}
