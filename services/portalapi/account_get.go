package portalapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"net/http"
)

type AccountGetRes struct {
	Account       *authenticator.Account  `json:"account"`
	Workspace     *entities.Workspace     `json:"workspace"`
	WorkspaceTier *entities.WorkspaceTier `json:"workspace_tier"`
}

// UseAccountGet
// @Tags		account
// @Router		/account/me			[get]
// @Success		200						{object}	AccountGetRes
// @Failure		default					{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseAccountGet() gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)
		ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
		wst := ctx.Value(authorizator.CtxWst).(*entities.WorkspaceTier)

		res := &AccountGetRes{Account: acc, Workspace: ws, WorkspaceTier: wst}
		ginctx.JSON(http.StatusOK, res)
	}
}
