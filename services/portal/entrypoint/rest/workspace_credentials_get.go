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

type WorkspaceCredentialsGetRes struct {
	*entities.WorkspaceCredentials
}

// UseWorkspaceCredentialsGet
// @Tags		workspace
// @Router		/workspace/me/credentials/{wsc_id}	[get]
// @Param		wsc_id								path		string						true	"credentials id"
// @Success		200									{object}	WorkspaceCredentialsGetRes
// @Failure		default								{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceCredentialsGet(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)

		id := ginctx.Param("wsc_id")
		ucreq := &usecase.WorkspaceCredentialsGetReq{
			WsId: ws.Id,
			Id:   id,
		}
		if err := ucreq.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := service.uc.WorkspaceCredentials().Get(ctx, ucreq)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCredentialsGetRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}