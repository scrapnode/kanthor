package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/internal/entities"
)

type WorkspaceGetRes struct {
	*Workspace
	Permissions []authorizator.Permission `json:"permissions"`
}

// UseWorkspaceGet
// @Tags		workspace
// @Router		/workspace/me			[get]
// @Success		200						{object}	WorkspaceGetRes
// @Failure		default					{object}	gateway.Error
// @Security	BearerAuth
// @Security	WsId
func UseWorkspaceGet(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		acc := ctx.Value(gateway.CtxAccount).(*authenticator.Account)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		permissions, err := service.infra.Authorizator.UserPermissionsInTenant(ws.Id, acc.Sub)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &WorkspaceGetRes{
			Workspace:   ToWorkspace(ws),
			Permissions: permissions,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
