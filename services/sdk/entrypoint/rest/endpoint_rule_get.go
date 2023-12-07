package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointRuleGetRes struct {
	*entities.EndpointRule
}

// UseEndpointRuleGet
// @Tags		endpoint rule
// @Router		/endpoint/{ep_id}/rule/{epr_id}	[get]
// @Param		ep_id							path		string					true	"endpoint id"
// @Param		epr_id							path		string					true	"rule id"
// @Success		200								{object}	EndpointRuleGetRes
// @Failure		default							{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointRuleGet(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		epId := ginctx.Param("ep_id")
		id := ginctx.Param("epr_id")
		in := &usecase.EndpointRuleGetIn{
			Ws:   ctx.Value(gateway.CtxWorkspace).(*entities.Workspace),
			EpId: epId,
			Id:   id,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.EndpointRule().Get(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleGetRes{out.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
