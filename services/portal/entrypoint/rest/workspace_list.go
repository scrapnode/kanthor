package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceListRes struct {
	Data  []Workspace `json:"data"`
	Count int64       `json:"count"`
} // @name WorkspaceListRes

// UseWorkspaceList
// @Tags			workspace
// @Router		/workspace		[get]
// @Success		200						{object}	WorkspaceGetRes
// @Failure		default				{object}	gateway.Err
// @Security	Authorization
func UseWorkspaceList(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)

		in := &usecase.WorkspaceListIn{
			AccId: acc.Sub,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Workspace().List(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &WorkspaceListRes{Data: make([]Workspace, 0), Count: out.Count}
		for _, ws := range out.Data {
			res.Data = append(res.Data, *ToWorkspace(&ws))
		}

		ginctx.JSON(http.StatusOK, res)
	}
}
