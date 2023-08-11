package portalapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
	"net/http"
)

type WorkspaceCredentialsExpireReq struct {
	Duration int64 `json:"duration"`
}

type WorkspaceCredentialsExpireRes struct {
	Id        string `json:"id"`
	ExpiredAt int64  `json:"expired_at"`
}

// UseWorkspaceCredentialsExpire
// @Tags		workspace
// @Router		/workspace/me/credentials/{wsc_id}/expiration	[put]
// @Param		wsc_id											path		string							true	"credentials id"
// @Param		props											body		WorkspaceCredentialsExpireReq	true	"credentials properties"
// @Success		200												{object}	WorkspaceCredentialsExpireRes
// @Failure		default											{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceCredentialsExpire(logger logging.Logger, validator validator.Validator, uc portaluc.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCredentialsExpireReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

		id := ginctx.Param("wsc_id")
		ucreq := &portaluc.WorkspaceCredentialsExpireReq{
			WorkspaceId: ws.Id,
			Id:          id,
			Duration:    req.Duration,
		}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.WorkspaceCredentials().Expire(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCredentialsExpireRes{
			Id:        ucres.Id,
			ExpiredAt: ucres.ExpiredAt,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}