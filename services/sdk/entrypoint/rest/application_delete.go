package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type ApplicationDeleteRes struct {
	*Application
} // @name ApplicationDeleteRes

// UseApplicationDelete
// @Tags		application
// @Router		/application/{app_id}	[delete]
// @Param		app_id					path		string					true	"application id"
// @Success		200						{object}	ApplicationDeleteRes
// @Failure		default					{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseApplicationDelete(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		id := ginctx.Param("app_id")
		in := &usecase.ApplicationDeleteIn{
			WsId: ws.Id,
			Id:   id,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Application().Delete(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &ApplicationDeleteRes{ToApplication(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
