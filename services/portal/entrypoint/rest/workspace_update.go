package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
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

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.WorkspaceUpdateIn{Id: ws.Id, Name: req.Name}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Workspace().Update(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceUpdateRes{out.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
