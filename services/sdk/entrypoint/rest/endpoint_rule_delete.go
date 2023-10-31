package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointRuleDeleteRes struct {
	*entities.EndpointRule
}

// UseEndpointRuleDelete
// @Tags		endpoint rule
// @Router		/endpoint/{ep_id}/rule/{epr_id}	[delete]
// @Param		ep_id							path		string					true	"endpoint id"
// @Param		epr_id							path		string					true	"rule id"
// @Success		200								{object}	EndpointRuleDeleteRes
// @Failure		default							{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointRuleDelete(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		epId := ginctx.Param("ep_id")
		id := ginctx.Param("epr_id")
		ucreq := &usecase.EndpointRuleDeleteReq{EpId: epId, Id: id}
		if err := ucreq.Validate(); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := service.uc.EndpointRule().Delete(ctx, ucreq)
		if err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleDeleteRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}