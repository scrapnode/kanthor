package portalapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"net/http"
)

type WorkspaceGetRes struct {
	Id       string `json:"id"`
	OwnerId  string `json:"owner_id"`
	Name     string `json:"name"`
	TierName string `json:"tier_name"`
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
		wst := ctx.Value(authorizator.CtxWst).(*entities.WorkspaceTier)

		res := &WorkspaceGetRes{
			Id:       ws.Id,
			OwnerId:  ws.OwnerId,
			Name:     ws.Name,
			TierName: wst.Name,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
