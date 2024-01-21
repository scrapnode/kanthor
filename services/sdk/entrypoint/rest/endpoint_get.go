package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointGetRes struct {
	*Endpoint
} // @name EndpointGetRes

// UseEndpointGet
// @Tags		endpoint
// @Router		/endpoint/{ep_id}	[get]
// @Param		ep_id				path		string					true	"endpoint id"
// @Success		200					{object}	EndpointGetRes
// @Failure		default				{object}	gateway.Err
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
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Endpoint().Get(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &EndpointGetRes{ToEndpoint(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
