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

type WorkspaceCreateReq struct {
	Name string `json:"name" default:"main"`
} // @name WorkspaceCreateReq

type WorkspaceCreateRes struct {
	*Workspace
} // @name WorkspaceCreateReq

// UseWorkspaceCreate
// @Tags		workspace
// @Router		/workspace			[post]
// @Param		props				body		WorkspaceCreateReq	true	"credentials properties"
// @Success		200					{object}	WorkspaceCreateRes
// @Failure		default				{object}	gateway.Error
// @Security	Authorization
func UseWorkspaceCreate(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		in := &usecase.WorkspaceCreateIn{
			AccId: acc.Sub,
			Name:  req.Name,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Workspace().Create(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceCreateRes{ToWorkspace(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
