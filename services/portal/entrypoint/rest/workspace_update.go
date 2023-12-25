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

type WorkspaceUpdateReq struct {
	Name string `json:"name" example:"another name"`
} // @name WorkspaceUpdateReq

type WorkspaceUpdateRes struct {
	*Workspace
} // @name WorkspaceUpdateRes

// UseWorkspaceUpdate
// @Tags		workspace
// @Router		/workspace/{ws_id}		[patch]
// @Param		ws_id					path		string						true	"workspace id"
// @Param		props					body		WorkspaceUpdateReq	true	"credentials properties"
// @Success		200						{object}	WorkspaceUpdateRes
// @Failure		default					{object}	gateway.Error
// @Security	Authorization
func UseWorkspaceUpdate(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		in := &usecase.WorkspaceUpdateIn{
			AccId: acc.Sub,
			Id:    ginctx.Param("ws_id"),
			Name:  req.Name,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Workspace().Update(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceUpdateRes{ToWorkspace(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
