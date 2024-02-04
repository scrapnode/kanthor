package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceGetRes struct {
	*Workspace
} // @name WorkspaceGetRes

// UseWorkspaceGet
// @Tags			workspace
// @Router		/workspace/{ws_id}	[get]
// @Param			ws_id								path			string						true	"workspace id"
// @Success		200									{object}	WorkspaceGetRes
// @Failure		default							{object}	gateway.Err
// @Security	Authorization
func UseWorkspaceGet(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		in := &usecase.WorkspaceGetIn{
			AccId: acc.Sub,
			Id:    ginctx.Param("ws_id"),
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Workspace().Get(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &WorkspaceGetRes{
			Workspace: ToWorkspace(out.Doc),
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
