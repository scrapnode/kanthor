package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	sdkuc "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

func UseAuth(uc sdkuc.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		credentials := ginctx.Request.Header.Get(authenticator.HeaderAuth)
		user, hash, err := authenticator.ParseBasicCredentials(credentials)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx := ginctx.MustGet("ctx").(context.Context)
		req := &sdkuc.WorkspaceAuthenticateReq{User: user, Hash: hash}
		res, err := uc.Workspace().Authenticate(ctx, req)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if res.Error != "" {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": res.Error})
			return
		}

		acc := &authenticator.Account{Sub: user, Name: res.WorkspaceCredentials.Name}
		ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)
		ctx = context.WithValue(ctx, authorizator.CtxWs, res.Workspace)
		ctx = context.WithValue(ctx, authorizator.CtxWst, res.WorkspaceTier)
		ginctx.Set("ctx", ctx)
		ginctx.Next()
	}
}
