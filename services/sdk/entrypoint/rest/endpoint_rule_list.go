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

type EndpointRuleListRes struct {
	Data  []EndpointRule
	Count int64
}

// UseEndpointRuleList
// @Tags		endpoint rule
// @Router		/rule			[get]
// @Param		ep_id			query		string					true	"endpoint id"
// @Param		_q				query		string					false	"search keyword"
// @Param		_limit			query		int						false	"limit returning records"	default(10)
// @Param		_page			query		int						false	"current requesting page"	default(0)
// @Success		200				{object}	EndpointRuleListRes
// @Failure		default			{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointRuleList(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var query gateway.Query
		if err := ginctx.BindQuery(&query); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("unable to parse your request query"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointRuleListIn{
			Query: &entities.Query{Search: query.Search, Page: query.Page, Limit: query.Limit},
			WsId:  ws.Id,
			EpId:  ginctx.Query("ep_id"),
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.EndpointRule().List(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleListRes{Data: make([]EndpointRule, 0), Count: out.Count}
		for _, epr := range out.Data {
			res.Data = append(res.Data, *ToEndpointRule(&epr))
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
