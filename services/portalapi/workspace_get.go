package portalapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
)

type WorkspaceGetRes struct {
	*entities.Workspace
}

// UseWorkspaceGet
// @Tags		workspace
// @Router		/workspace/me			[get]
// @Success		200						{object}	WorkspaceGetRes
// @Failure		default					{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceGet() gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

		res := &WorkspaceGetRes{Workspace: ws}
		ginctx.JSON(http.StatusOK, res)
	}
}
