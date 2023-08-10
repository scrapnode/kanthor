package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
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

		ctx := ginctx.MustGet("ctx").(context.Context)
		ctx = context.WithValue(ctx, authenticator.CtxAcc, acc)
		ginctx.Set("ctx", ctx)
		ginctx.Next()
	}
}
