package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	sdkuc "github.com/scrapnode/kanthor/usecases/sdk"
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
	uc sdkuc.Sdk,
) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)

		// authenticate sdk
		credentials := ginctx.Request.Header.Get(authenticator.HeaderAuth)
		ctx, err := sdkauth(uc, ctx, credentials)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unable to verify your credentials"})
			return
		}

		// @TODO: authenticate portal

		ginctx.Set(gateway.KeyContext, ctx)
		ginctx.Next()
	}
}

func sdkauth(
	uc sdkuc.Sdk,
	ctx context.Context,
	credentials string,
) (context.Context, error) {
	user, pass, err := authenticator.ParseBasicCredentials(credentials)
	if err != nil {
		return ctx, err
	}

	ucreq := &sdkuc.WorkspaceCredentialsAuthenticateReq{User: user, Pass: pass}
	if err := ucreq.Validate(); err != nil {
		return ctx, err
	}

	ucres, err := uc.WorkspaceCredentials().Authenticate(ctx, ucreq)
	if err != nil {
		return ctx, err
	}

	acc := &authenticator.Account{Sub: user, Name: ucres.WorkspaceCredentials.Name}
	ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)
	ctx = context.WithValue(ctx, gateway.CtxWs, ucres.Workspace)
	return ctx, nil
}
