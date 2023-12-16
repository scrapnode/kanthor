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

type EndpointRuleUpdateReq struct {
	Name string `json:"name" binding:"required"`
}

type EndpointRuleUpdateRes struct {
	*entities.EndpointRule
}

// UseEndpointRuleUpdate
// @Tags		endpoint rule
// @Router		/endpoint/{ep_id}/rule/{epr_id}	[put]
// @Param		ep_id							path		string					true	"endpoint id"
// @Param		epr_id							path		string					true	"rule id"
// @Param		props							body		EndpointRuleUpdateReq	true	"rule properties"
// @Success		200								{object}	EndpointRuleUpdateRes
// @Failure		default							{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointRuleUpdate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointRuleUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		epId := ginctx.Param("ep_id")
		id := ginctx.Param("epr_id")
		in := &usecase.EndpointRuleUpdateIn{
			Ws:   ctx.Value(gateway.CtxWorkspace).(*entities.Workspace),
			EpId: epId,
			Id:   id,
			Name: req.Name,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.EndpointRule().Update(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleUpdateRes{out.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
