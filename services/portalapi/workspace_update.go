package portalapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
)

type WorkspaceUpdateReq struct {
	Name string `json:"name" binding:"required"`
}

type WorkspaceUpdateRes struct {
	*entities.Workspace
}

// UseWorkspaceUpdate
// @Tags		workspace
// @Router		/workspace/me			[put]
// @Param		props					body		WorkspaceUpdateReq	true	"credentials properties"
// @Success		200						{object}	WorkspaceUpdateRes
// @Failure		default					{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceUpdate(logger logging.Logger, uc portaluc.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

		ucreq := &portaluc.WorkspaceUpdateReq{Id: ws.Id, Name: req.Name}
		if err := ucreq.Validate(); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Workspace().Update(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceUpdateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
