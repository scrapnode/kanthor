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

type EndpointRuleCreateReq struct {
	Name string `json:"name" binding:"required"`

	Priority            int32  `json:"priority"`
	Exclusionary        bool   `json:"exclusionary"`
	ConditionSource     string `json:"condition_source" binding:"required"`
	ConditionExpression string `json:"condition_expression" binding:"required"`
}

type EndpointRuleCreateRes struct {
	*entities.EndpointRule
}

// UseEndpointRuleCreate
// @Tags		endpoint rule
// @Router		/endpoint/{ep_id}/rule	[post]
// @Param		ep_id					path		string					true	"endpoint id"
// @Param		props					body		EndpointRuleCreateReq	true	"rule properties"
// @Success		201						{object}	EndpointRuleCreateRes
// @Failure		default					{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointRuleCreate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointRuleCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		epId := ginctx.Param("ep_id")

		in := &usecase.EndpointRuleCreateIn{
			Ws:                  ctx.Value(gateway.CtxWorkspace).(*entities.Workspace),
			EpId:                epId,
			Name:                req.Name,
			Priority:            req.Priority,
			Exclusionary:        req.Exclusionary,
			ConditionSource:     req.ConditionSource,
			ConditionExpression: req.ConditionExpression,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.EndpointRule().Create(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointRuleCreateRes{out.Doc}
		ginctx.JSON(http.StatusCreated, res)
	}
}
