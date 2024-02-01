package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/internal/entities"
)

func UseWorkspace(resolve func(ctx context.Context, acc *authenticator.Account, id string) (*entities.Workspace, error)) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		var wsId string
		if id, has := acc.Metadata[gateway.MetaWorkspaceId]; has && id != "" {
			wsId = id
		}
		if id := ginctx.Request.Header.Get(authenticator.HeaderAuthnWorkspace); id != "" {
			wsId = id
		}

		if wsId == "" {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gateway.ErrorString("GATEWAY.MIDDLEWARE.WORKSPACE.UNKNOWN.ERROR"))
			return
		}

		workspace, err := resolve(ctx, acc, wsId)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ginctx.Set(gateway.Ctx, context.WithValue(ctx, gateway.CtxWorkspace, workspace))
		ginctx.Next()
	}
}
