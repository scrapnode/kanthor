package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
)

func UseAuth(auth authenticator.Authenticator) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		// authorization format is <auth-scheme> <authorization-parameters>
		scheme, credentials := authenticator.Parse(ginctx.Request.Header.Get(authenticator.HeaderAuth))

		// use custom scheme instead of the standard at https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml
		if custom := ginctx.Request.Header.Get(authenticator.HeaderAuthScheme); custom != "" {
			scheme = custom
		}
		request := &authenticator.Request{Credentials: credentials, Metadata: map[string]string{}}
		for key, value := range ginctx.Request.Header {
			request.Metadata[key] = value[0]
		}
		acc, err := auth.Authenticate(ctx, scheme, request)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ginctx.Set(gateway.Ctx, context.WithValue(ctx, gateway.CtxAccount, acc))
		ginctx.Next()
	}
}
