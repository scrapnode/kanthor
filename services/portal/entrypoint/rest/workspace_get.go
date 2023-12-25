package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceGetRes struct {
	*Workspace
	Permissions []authorizator.Permission `json:"permissions"`
} // @name WorkspaceGetRes

// UseWorkspaceGet
// @Tags		workspace
// @Router		/workspace/{ws_id}	[get]
// @Param		ws_id				path		string						true	"workspace id"
// @Success		200					{object}	WorkspaceGetRes
// @Failure		default				{object}	gateway.Error
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
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Workspace().Get(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		permissions, err := service.infra.Authorizator.UserPermissionsInTenant(out.Doc.Id, acc.Sub)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceGetRes{
			Workspace:   ToWorkspace(out.Doc),
			Permissions: permissions,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
