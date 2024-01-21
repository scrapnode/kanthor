package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceCreateReq struct {
	Name string `json:"name" default:"main"`
} // @name WorkspaceCreateReq

type WorkspaceCreateRes struct {
	*Workspace
} // @name WorkspaceCreateRes

// UseWorkspaceCreate
// @Tags		workspace
// @Router		/workspace			[post]
// @Param		payload				body		WorkspaceCreateReq	true	"workspace payload"
// @Success		200					{object}	WorkspaceCreateRes
// @Failure		default				{object}	gateway.Err
// @Security	Authorization
func UseWorkspaceCreate(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		in := &usecase.WorkspaceCreateIn{
			AccId: acc.Sub,
			Name:  req.Name,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Workspace().Create(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &WorkspaceCreateRes{ToWorkspace(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
