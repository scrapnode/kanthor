package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type AccountSetupReq struct {
	WorkspaceName string `json:"workspace_name"`
}

type AccountSetupRes struct {
	Account   *authenticator.Account `json:"account"`
	Workspace *entities.Workspace    `json:"workspace"`
}

// UseAccountSetup
// @Tags		account
// @Router		/account/setup			[post]
// @Param		props				body		AccountSetupReq	true	"setup options"
// @Success		200					{object}	AccountSetupRes
// @Failure		default				{object}	gateway.Error
// @Security	BearerAuth
func UseAccountSetup(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req AccountSetupReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		in := &usecase.AccountSetupIn{AccountId: acc.Sub, WorkspaceName: req.WorkspaceName}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Account().Setup(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &AccountSetupRes{Account: acc, Workspace: out.Workspace}
		ginctx.JSON(http.StatusOK, res)
	}
}
