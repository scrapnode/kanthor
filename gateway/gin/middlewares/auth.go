package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
)

func UseAuth(auth authenticator.Authenticator, defaultEngine string) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		credentials := ginctx.Request.Header.Get(authenticator.HeaderAuthCredentials)
		request := &authenticator.Request{Credentials: credentials, Metadata: map[string]string{}}
		for key, value := range ginctx.Request.Header {
			request.Metadata[key] = value[0]
		}

		engine := defaultEngine
		if selectEngine := ginctx.Request.Header.Get(authenticator.HeaderAuthEngine); selectEngine != "" {
			engine = selectEngine
		}
		acc, err := auth.Authenticate(engine, ctx, request)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ginctx.Set(gateway.Ctx, context.WithValue(ctx, gateway.CtxAccount, acc))
		ginctx.Next()
	}
}
