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

type EndpointGetRes struct {
	*Endpoint
}

// UseEndpointGet
// @Tags		endpoint
// @Router		/endpoint/{ep_id}	[get]
// @Param		app_id				query		string					true	"application id"
// @Param		ep_id				path		string					true	"endpoint id"
// @Success		200					{object}	EndpointGetRes
// @Failure		default				{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointGet(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointGetIn{
			WsId: ws.Id,
			Id:   ginctx.Param("ep_id"),
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

		res := &EndpointGetRes{ToEndpoint(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
