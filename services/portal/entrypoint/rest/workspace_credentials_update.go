package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceCredentialsUpdateReq struct {
	Name string `json:"name" binding:"required"`
}

type WorkspaceCredentialsUpdateRes struct {
	*entities.WorkspaceCredentials
}

// UseWorkspaceCredentialsUpdate
// @Tags		workspace
// @Router		/workspace/me/credentials/{wsc_id}	[put]
// @Param		wsc_id								path		string							true	"credentials id"
// @Param		props								body		WorkspaceCredentialsUpdateReq	true	"credentials properties"
// @Success		200									{object}	WorkspaceCredentialsUpdateRes
// @Failure		default								{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceCredentialsUpdate(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCredentialsUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

		id := ginctx.Param("wsc_id")
		in := &usecase.WorkspaceCredentialsUpdateIn{
			WsId: ws.Id,
			Id:   id,
			Name: req.Name,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.WorkspaceCredentials().Update(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCredentialsUpdateRes{out.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
