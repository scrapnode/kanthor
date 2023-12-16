package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type ApplicationGetRes struct {
	*entities.Application
}

// UseApplicationGet
// @Tags		application
// @Router		/application/{app_id}	[get]
// @Param		app_id					path		string					true	"application id"
// @Success		200						{object}	ApplicationGetRes
// @Failure		default					{object}	gateway.Error
// @Security	BasicAuth
func UseApplicationGet(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		id := ginctx.Param("app_id")
		in := &usecase.ApplicationGetIn{
			WsId: ws.Id,
			Id:   id,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Application().Get(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &ApplicationGetRes{out.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
