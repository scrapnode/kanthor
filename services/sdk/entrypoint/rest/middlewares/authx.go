package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/authenticator"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

var (
	HeaderAuthScheme  = "kanthor-auth-scheme"
	AuthSchemeForward = "forward"
)

type PortalAuth struct {
	Error     string                 `json:"error"`
	Account   *authenticator.Account `json:"account"`
	Workspace *entities.Workspace    `json:"workspace"`
}

func UseAuth(
	authz authorizator.Authorizator,
	uc usecase.Sdk,
) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)

		// authenticate sdk
		credentials := ginctx.Request.Header.Get(authenticator.HeaderAuth)
		ctx, err := auth(uc, ctx, credentials)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ginctx.Set(gateway.KeyContext, ctx)
		ginctx.Next()
	}
}

func auth(
	uc usecase.Sdk,
	ctx context.Context,
	credentials string,
) (context.Context, error) {
	user, pass, err := authenticator.ParseBasicCredentials(credentials)
	if err != nil {
		return ctx, err
	}

	in := &usecase.WorkspaceCredentialsAuthenticateIn{User: user, Pass: pass}
	if err := in.Validate(); err != nil {
		return ctx, err
	}

	out, err := uc.WorkspaceCredentials().Authenticate(ctx, in)
	if err != nil {
		return ctx, err
	}

	acc := &authenticator.Account{Sub: user, Name: out.WorkspaceCredentials.Name}
	ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)
	ctx = context.WithValue(ctx, gateway.CtxWs, out.Workspace)
	return ctx, nil
}
