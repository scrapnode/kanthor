package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
)

type AccountGetRes struct {
	Account *authenticator.Account `json:"account"`
} // @name AccountGetRes

// UseAccountGet
// @Tags		account
// @Router		/account/me			[get]
// @Success		200					{object}	AccountGetRes
// @Failure		default				{object}	gateway.Error
// @Security	Authorization
func UseAccountGet(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		res := &AccountGetRes{Account: acc}
		ginctx.JSON(http.StatusOK, res)
	}
}
