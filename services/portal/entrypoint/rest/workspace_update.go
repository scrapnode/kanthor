package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
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
func UseWorkspaceUpdate(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

		ucreq := &usecase.WorkspaceUpdateReq{Id: ws.Id, Name: req.Name}
		if err := ucreq.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := service.uc.Workspace().Update(ctx, ucreq)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceUpdateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}