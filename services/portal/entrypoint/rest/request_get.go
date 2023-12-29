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

type RequestGetRes struct {
	*Request
} // @name RequestGetRes

// UseRequestGet
// @Tags		request
// @Router		/request/{id}		[get]
// @Param		ep_id				query		string			true	"epndpoint id"
// @Param		msg_id				query		string			true	"msg id"
// @Success		200					{object}	RequestGetRes
// @Failure		default				{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseRequestGet(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.RequestGetIn{
			WsId:  ws.Id,
			EpId:  ginctx.Query("ep_id"),
			MsgId: ginctx.Query("msg_id"),
			Id:    ginctx.Param("req_id"),
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Request().Get(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &RequestGetRes{ToRequest(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
