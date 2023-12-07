package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/internal/domain/entities"
)

func UseWorkspace(resolve func(ctx context.Context, id string) (*entities.Workspace, error)) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		id := ginctx.Request.Header.Get(authorizator.HeaderAuthWorkspace)
		workspace, err := resolve(ctx, id)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ginctx.Set(gateway.Ctx, context.WithValue(ctx, gateway.CtxWorkspace, workspace))
		ginctx.Next()
	}
}
