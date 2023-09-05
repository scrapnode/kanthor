package portalapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"net/http"
)

type AccountGetRes struct {
	Account *authenticator.Account `json:"account"`
}

// UseAccountGet
// @Tags		account
// @Router		/account/me			[get]
// @Success		200					{object}	AccountGetRes
// @Failure		default				{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseAccountGet() gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)

		res := &AccountGetRes{Account: acc}
		ginctx.JSON(http.StatusOK, res)
	}
}
