package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointUpdateReq struct {
	Name   string `json:"name"`
	Method string `json:"method"`
} // @name EndpointUpdateReq

type EndpointUpdateRes struct {
	*Endpoint
} // @name EndpointUpdateRes

// UseEndpointUpdate
// @Tags		endpoint
// @Router		/endpoint/{ep_id}	[patch]
// @Param		ep_id				path		string					true	"endpoint id"
// @Param		payload				body		EndpointUpdateReq		true	"endpoint payload"
// @Success		200					{object}	EndpointUpdateRes
// @Failure		default				{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointUpdate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointUpdateIn{
			WsId:   ws.Id,
			Id:     ginctx.Param("ep_id"),
			Name:   req.Name,
			Method: req.Method,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Endpoint().Update(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &EndpointUpdateRes{ToEndpoint(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
