package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/internal/domain/entities"
)

func UseAuthz(authz authorizator.Authorizator) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		acc, ok := ctx.Value(gateway.CtxAccount).(*authenticator.Account)
		if !ok {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown account"})
			return
		}

		workspace, ok := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)
		if !ok {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown workspace"})
			return
		}

		obj := ginctx.FullPath()
		act := ginctx.Request.Method
		ok, err := authz.Enforce(workspace.Id, acc.Sub, obj, act)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !ok {
			ginctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denined"})
			return
		}

		ginctx.Set(gateway.Ctx, ctx)
		ginctx.Next()
	}
}
