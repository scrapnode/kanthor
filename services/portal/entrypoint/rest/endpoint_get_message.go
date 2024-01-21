package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type EndpointGetMessageRes struct {
	*EndpointMessage
} // @name EndpointGetMessageRes

// UseEndpointGetMessage
// @Tags		endpoint
// @Router		/endpoint/{ep_id}/message/{msg_id}			[get]
// @Success		200											{object}	EndpointGetMessageRes
// @Failure		default										{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointGetMessage(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointGetMessageIn{
			WsId:  ws.Id,
			EpId:  ginctx.Param("ep_id"),
			MsgId: ginctx.Param("msg_id"),
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Endpoint().GetMessage(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &EndpointGetMessageRes{ToEndpointMessage(out.Doc, out.Requests, out.Responses)}
		ginctx.JSON(http.StatusOK, res)
	}
}
