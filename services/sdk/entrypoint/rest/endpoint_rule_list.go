package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/structure"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointRuleListRes struct {
	*structure.ListRes[entities.EndpointRule]
}

// UseEndpointRuleList
// @Tags		endpoint rule
// @Router		/endpoint/{ep_id}/rule	[get]
// @Param		ep_id					path		string					true	"endpoint id"
// @Param		_cursor					query		string					false	"current query cursor"					minlength(29) maxlength(32)
// @Param		_q						query		string					false	"search keyword" 						minlength(2)  maxlength(32)
// @Param		_limit					query		int						false	"limit returning records"				minimum(5)    maximum(30)
// @Param		_id						query		[]string				false	"only return records with selected ids"
// @Success		200						{object}	EndpointRuleListRes
// @Failure		default					{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointRuleList(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		epId := ginctx.Param("ep_id")

		in := &usecase.EndpointRuleListIn{
			Ws:      ctx.Value(gateway.CtxWorkspace).(*entities.Workspace),
			EpId:    epId,
			ListReq: ginctx.MustGet("list_req").(*structure.ListReq),
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

		res := &EndpointRuleListRes{ListRes: out.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}
