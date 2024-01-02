package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type EndpointListMessageRes struct {
	Data []EndpointMessage `json:"data"`
} // @name EndpointListMessageRes

// UseEndpointListMessage
// @Tags		endpoint
// @Router		/endpoint/{ep_id}/message			[get]
// @Param		_limit								query		int				false	"limit returning records" 	default(10)
// @Param		_start								query		int64			false	"starting time to scan in milliseconds" example(1669914060000)
// @Param		_end								query		int64			false	"ending time to scan in milliseconds" 	example(1985533260000)
// @Success		200									{object}	EndpointListMessageRes
// @Failure		default								{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointListMessage(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req gateway.Query
		if err := ginctx.BindQuery(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("unable to parse your request query"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)
		in := &usecase.EndpointListMessageIn{
			ScanningQuery: entities.ScanningQueryFromGatewayQuery(&req, service.infra.Timer),
			WsId:          ws.Id,
			EpId:          ginctx.Param("ep_id"),
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Endpoint().ListMessage(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointListMessageRes{Data: make([]EndpointMessage, 0)}
		for _, doc := range out.Data {
			res.Data = append(res.Data, *ToEndpointMessage(&doc, nil, nil))
		}
		ginctx.JSON(http.StatusOK, res)
	}
}