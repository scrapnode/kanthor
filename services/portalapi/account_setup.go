package portalapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validation"
	"github.com/scrapnode/kanthor/pkg/utils"
	portaluc "github.com/scrapnode/kanthor/usecases/portal"
)

type AccountSetupReq struct {
}

type AccountSetupRes struct {
	Account       *authenticator.Account  `json:"account"`
	Workspace     *entities.Workspace     `json:"workspace"`
	WorkspaceTier *entities.WorkspaceTier `json:"workspace_tier"`
}

// UseAccountSetup
// @Tags		account
// @Router		/account/me			[put]
// @Param		props				body		AccountSetupReq	true	"setup options"
// @Success		200					{object}	AccountSetupRes
// @Failure		default				{object}	gateway.Error
// @Security	BearerAuth
func UseAccountSetup(logger logging.Logger, validator validation.Validator, uc portaluc.Portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		acc := ctx.Value(authenticator.CtxAcc).(*authenticator.Account)

		ucreq := &portaluc.AccountSetupReq{AccountId: acc.Sub}
		if err := validator.Struct(ucreq); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Account().Setup(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &AccountSetupRes{Account: acc, Workspace: ucres.Workspace, WorkspaceTier: ucres.WorkspaceTier}
		ginctx.JSON(http.StatusOK, res)
	}
}
