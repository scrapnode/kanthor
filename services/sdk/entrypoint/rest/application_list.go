package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type ApplicationListRes struct {
	Data  []Application `json:"data"`
	Count int64         `json:"count"`
} // @name ApplicationListRes

// UseApplicationList
// @Tags			application
// @Router		/application	[get]
// @Param			id						query			[]string					false	"list by ids"
// @Param			_q						query			string						false	"search keyword"
// @Param			_limit				query			int								false	"limit returning records"	default(10)
// @Param			_page					query			int								false	"current requesting page"	default(0)
// @Success		200						{object}	ApplicationListRes
// @Failure		default				{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseApplicationList(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var query gateway.Query
		if err := ginctx.BindQuery(&query); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.ApplicationListIn{
			PagingQuery: entities.PagingQueryFromGatewayQuery(&query),
			WsId:        ws.Id,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Application().List(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &ApplicationListRes{Data: make([]Application, 0), Count: out.Count}
		for _, app := range out.Data {
			res.Data = append(res.Data, *ToApplication(&app))
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
