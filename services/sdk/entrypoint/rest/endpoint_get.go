package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointGetRes struct {
	*entities.Endpoint
}

// UseEndpointGet
// @Tags		endpoint
// @Router		/application/{app_id}/endpoint/{ep_id}	[get]
// @Param		app_id									path		string					true	"application id"
// @Param		ep_id									path		string					true	"endpoint id"
// @Success		200										{object}	EndpointGetRes
// @Failure		default									{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointGet(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		appId := ginctx.Param("app_id")
		id := ginctx.Param("ep_id")
		in := &usecase.EndpointGetIn{
			WsId:  ws.Id,
			AppId: appId,
			Id:    id,
		}
		if err := in.Validate(); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Endpoint().Get(ctx, in)
		if err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointGetRes{out.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
