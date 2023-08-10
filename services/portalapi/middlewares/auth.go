package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"net/http"
)

func UseAuth(engine authenticator.Authenticator) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		credentials := ginctx.Request.Header.Get(authenticator.HeaderAuth)
		acc, err := engine.Verify(credentials)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)
		ginctx.Set(gateway.KeyCtx, ctx)
		ginctx.Next()
	}
}
