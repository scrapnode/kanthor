package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/domain/entities"
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
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		res := &WorkspaceGetRes{Workspace: ws}
		ginctx.JSON(http.StatusOK, res)
	}
}
