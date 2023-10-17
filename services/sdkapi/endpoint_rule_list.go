package sdkapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
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
func UseEndpointRuleList(logger logging.Logger, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		epId := ginctx.Param("ep_id")

		ucreq := &usecase.EndpointRuleListReq{
			EpId:    epId,
			ListReq: ginctx.MustGet("list_req").(*structure.ListReq),
		}
		if err := ucreq.Validate(); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.EndpointRule().List(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleListRes{ListRes: ucres.ListRes}
		ginctx.JSON(http.StatusOK, res)
	}
}
