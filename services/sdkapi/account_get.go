package sdkapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
)

type AccountGetReq struct {
	WorkspcaeId string `form:"workspace_id" binding:"required,startswith=ws_"`
}

type AccountGetRes struct {
	Account     *authenticator.Account    `json:"account"`
	Permissions []authorizator.Permission `json:"permissions"`
}

// UseAccountGet
// @Tags		account
// @Router		/account/me			[get]
// @Success		200					{object}	AccountGetRes
// @Failure		default				{object}	gateway.Error
// @Security	BearerAuth
func UseAccountGet(service *sdkapi) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)

		var req AccountGetReq
		if err := ginctx.BindQuery(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("unable to parse your request query"))
			return
		}

		permissions, err := service.infra.Authorizator.UserPermissionsInTenant(req.WorkspcaeId, acc.Sub)
		if err != nil {
			service.logger.Errorw(err.Error(), "worksapce_id", req.WorkspcaeId)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("unable to get your permissions"))
			return
		}

		res := &AccountGetRes{Account: acc, Permissions: permissions}
		ginctx.JSON(http.StatusOK, res)
	}
}
