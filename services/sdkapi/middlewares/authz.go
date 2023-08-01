package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	sdkuc "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

func UseAuthz(authz authorizator.Authorizator) gin.HandlerFunc {
	return func(ginctx *gin.Context) {

		ctx := ginctx.MustGet("ctx").(context.Context)
		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)
		ws := ctx.Value(sdkuc.CtxWs).(*entities.Workspace)
		obj := ginctx.FullPath() // form of /application/:app_id
		act := ginctx.Request.Method
		ok, err := authz.Enforce(acc.Sub, ws.Id, obj, act)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you have no permission to perform the action"})
			return
		}

		ginctx.Next()
	}
}
