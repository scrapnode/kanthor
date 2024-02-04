package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceUpdateReq struct {
	Name string `json:"name" example:"another name"`
} // @name WorkspaceUpdateReq

type WorkspaceUpdateRes struct {
	*Workspace
} // @name WorkspaceUpdateRes

// UseWorkspaceUpdate
// @Tags			workspace
// @Router		/workspace/{ws_id}		[patch]
// @Param			ws_id									path			string							true	"workspace id"
// @Param			payload								body			WorkspaceUpdateReq	true	"workspace payload"
// @Success		200										{object}	WorkspaceUpdateRes
// @Failure		default								{object}	gateway.Err
// @Security	Authorization
func UseWorkspaceUpdate(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req WorkspaceUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
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
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Workspace().Update(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &WorkspaceUpdateRes{ToWorkspace(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
