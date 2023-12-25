package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type AccountGetRes struct {
	Account    *authenticator.Account `json:"account"`
	Workspaces []Workspace            `json:"workspaces"`
} // @name AccountGetRes

// UseAccountGet
// @Tags		account
// @Router		/account			[get]
// @Success		200					{object}	AccountGetRes
// @Failure		default				{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseAccountGet(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)
		res := &AccountGetRes{Account: acc, Workspaces: make([]Workspace, 0)}

		// we don't need to validate input here because we know it must be valid
		in := &usecase.WorkspaceListIn{AccId: acc.Sub}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Workspace().List(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}
		for _, ws := range out.Data {
			res.Workspaces = append(res.Workspaces, *ToWorkspace(&ws))
		}

		ginctx.JSON(http.StatusOK, res)
	}
}
