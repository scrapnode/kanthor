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

type WorkspaceCredentialsUpdateReq struct {
	Name      string `json:"name" default:"swagger demo update"`
	ExpiredAt int64  `json:"expired_at" default:"1893456000000"`
} // @name WorkspaceCredentialsUpdateReq

type WorkspaceCredentialsUpdateRes struct {
	*WorkspaceCredentials
} // @name WorkspaceCredentialsUpdateRes

// UseWorkspaceCredentialsUpdate
// @Tags		credentials
// @Router		/credentials/{wsc_id}	[patch]
// @Param		wsc_id					path		string							true	"credentials id"
// @Param		payload					body		WorkspaceCredentialsUpdateReq	true	"credentials payload"
// @Success		200						{object}	WorkspaceCredentialsUpdateRes
// @Failure		default					{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseWorkspaceCredentialsUpdate(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCredentialsUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.WorkspaceCredentialsUpdateIn{
			WsId:      ws.Id,
			Id:        ginctx.Param("wsc_id"),
			Name:      req.Name,
			ExpiredAt: req.ExpiredAt,
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

		res := &WorkspaceCredentialsUpdateRes{ToWorkspaceCredentials(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
