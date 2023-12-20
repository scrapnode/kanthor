package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/internal/entities"
)

func UseWorkspace(resolve func(ctx context.Context, id string) (*entities.Workspace, error)) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		id := ginctx.Request.Header.Get(authorizator.HeaderAuthWorkspace)
		// fallback to default workspace id of request account
		if ws, has := acc.Metadata[string(gateway.CtxWorkspaceId)]; has && id == "" {
			id = ws
		}
		if id == "" {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown selected workspace"})
			return
		}

		workspace, err := resolve(ctx, id)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ginctx.Set(gateway.Ctx, context.WithValue(ctx, gateway.CtxWorkspace, workspace))
		ginctx.Next()
	}
}
