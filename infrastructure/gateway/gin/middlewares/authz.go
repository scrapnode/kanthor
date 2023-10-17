package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
)

func UseAuthz(authz authorizator.Authorizator) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)

		// already authorized and passed, go to next routine
		if ok, has := ctx.Value(gateway.CtxAuhzOk).(bool); has && ok {
			ginctx.Next()
			return
		}

		// starting authorize process
		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)
		ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

		obj := ginctx.FullPath() // form of /application/:app_id
		act := ginctx.Request.Method
		ok, err := authz.Enforce(ws.Id, acc.Sub, obj, act)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you have no permission to perform the action"})
			return
		}

		ctx = context.WithValue(ctx, gateway.CtxAuhzOk, true)
		ginctx.Set(gateway.KeyContext, ctx)
		ginctx.Next()
	}
}
