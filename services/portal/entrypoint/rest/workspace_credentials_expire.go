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

type WorkspaceCredentialsExpireReq struct {
	Duration int64 `json:"duration" default:"1800000"`
} // @name WorkspaceCredentialsExpireReq

type WorkspaceCredentialsExpireRes struct {
	Id        string `json:"id"`
	ExpiredAt int64  `json:"expired_at"`
} // @name WorkspaceCredentialsExpireRes

// UseWorkspaceCredentialsExpire
// @Tags		credentials
// @Router		/credentials/{wsc_id}/expiration	[put]
// @Param		wsc_id								path		string							true	"credentials id"
// @Param		props								body		WorkspaceCredentialsExpireReq	true	"credentials properties"
// @Success		200									{object}	WorkspaceCredentialsExpireRes
// @Failure		default								{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseWorkspaceCredentialsExpire(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCredentialsExpireReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		id := ginctx.Param("wsc_id")
		in := &usecase.WorkspaceCredentialsExpireIn{
			WsId:     ws.Id,
			Id:       id,
			Duration: req.Duration,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.WorkspaceCredentials().Expire(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCredentialsExpireRes{
			Id:        out.Id,
			ExpiredAt: out.ExpiredAt,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
