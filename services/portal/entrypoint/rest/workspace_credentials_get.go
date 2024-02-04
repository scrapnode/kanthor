package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type WorkspaceCredentialsGetRes struct {
	*WorkspaceCredentials
} // @name WorkspaceCredentialsGetRes

// UseWorkspaceCredentialsGet
// @Tags			credentials
// @Router		/credentials/{wsc_id}	[get]
// @Param			wsc_id								path			string										true	"credentials id"
// @Success		200										{object}	WorkspaceCredentialsGetRes
// @Failure		default								{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseWorkspaceCredentialsGet(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.WorkspaceCredentialsGetIn{
			WsId: ws.Id,
			Id:   ginctx.Param("wsc_id"),
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.WorkspaceCredentials().Get(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &WorkspaceCredentialsGetRes{ToWorkspaceCredentials(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
