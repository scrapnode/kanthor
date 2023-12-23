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

type EndpointUpdateReq struct {
	Name   string `json:"name"`
	Method string `json:"method"`
}

type EndpointUpdateRes struct {
	*Endpoint
}

// UseEndpointUpdate
// @Tags		endpoint
// @Router		/endpoint/{ep_id}	[patch]
// @Param		ep_id				path		string					true	"endpoint id"
// @Param		props				body		EndpointUpdateReq		true	"endpoint properties"
// @Success		200					{object}	EndpointUpdateRes
// @Failure		default				{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointUpdate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
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
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Endpoint().Update(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointUpdateRes{ToEndpoint(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
