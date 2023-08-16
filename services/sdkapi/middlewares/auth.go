package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	sdkuc "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

func UseAuth(validator validator.Validator, uc sdkuc.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)

		credentials := ginctx.Request.Header.Get(authenticator.HeaderAuth)
		ctx, err := basic(validator, uc, ctx, credentials)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ginctx.Set(gateway.KeyCtx, ctx)
		ginctx.Next()
	}
}

func basic(
	validator validator.Validator,
	uc sdkuc.Sdk,
	ctx context.Context,
	credentials string,
) (context.Context, error) {
	user, hash, err := authenticator.ParseBasicCredentials(credentials)
	if err != nil {
		return ctx, err
	}

	req := &sdkuc.WorkspaceCredentialsAuthenticateReq{User: user, Hash: hash}
	if err := validator.Struct(req); err != nil {
		return ctx, err
	}

	res, err := uc.WorkspaceCredentials().Authenticate(ctx, req)
	if err != nil {
		return ctx, err
	}

	acc := &authenticator.Account{Sub: user, Name: res.WorkspaceCredentials.Name}
	ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)
	ctx = context.WithValue(ctx, authorizator.CtxWs, res.Workspace)
	ctx = context.WithValue(ctx, authorizator.CtxWst, res.WorkspaceTier)
	return ctx, nil
}
