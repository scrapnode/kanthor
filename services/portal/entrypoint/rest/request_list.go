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

type RequestListReq struct {
	*gateway.Query
	EpId string `json:"ep_id" form:"ep_id"`
} // @name RequestListReq

type RequestListRes struct {
	Data []Request `json:"data"`
} // @name RequestListRes

// UseRequestList
// @Tags		request
// @Router		/request		[get]
// @Param		ep_id			query		string			true	"epndpoint id"
// @Param		_limit			query		int				false	"limit returning records" 	default(10)
// @Param		_start			query		int64			false	"starting time to scan in milliseconds" example(1669914060000)
// @Param		_end			query		int64			false	"ending time to scan in milliseconds" 	example(1985533260000)
// @Success		200				{object}	RequestListRes
// @Failure		default			{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseRequestList(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req RequestListReq
		if err := ginctx.BindQuery(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("unable to parse your request query"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)

		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)
		in := &usecase.RequestListIn{
			ScanningQuery: entities.ScanningQueryFromGatewayQuery(req.Query, service.infra.Timer),
			WsId:          ws.Id,
			EpId:          req.EpId,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Request().List(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &RequestListRes{Data: make([]Request, 0)}
		for _, ws := range out.Data {
			res.Data = append(res.Data, *ToRequest(&ws))
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
