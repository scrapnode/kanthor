package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointRuleUpdateReq struct {
	Name string `json:"name"`

	Priority            int32  `json:"priority"`
	Exclusionary        bool   `json:"exclusionary"`
	ConditionSource     string `json:"condition_source"`
	ConditionExpression string `json:"condition_expression"`
} // @name EndpointRuleUpdateReq

type EndpointRuleUpdateRes struct {
	*EndpointRule
} // @name EndpointRuleUpdateRes

// UseEndpointRuleUpdate
// @Tags		endpoint rule
// @Router		/rule/{epr_id}	[patch]
// @Param		epr_id			path		string					true	"rule id"
// @Param		payload			body		EndpointRuleUpdateReq	true	"rule payload"
// @Success		200				{object}	EndpointRuleUpdateRes
// @Failure		default			{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointRuleUpdate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointRuleUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointRuleUpdateIn{
			WsId:                ws.Id,
			Id:                  ginctx.Param("epr_id"),
			Name:                req.Name,
			Priority:            req.Priority,
			Exclusionary:        req.Exclusionary,
			ConditionSource:     req.ConditionSource,
			ConditionExpression: req.ConditionExpression,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.EndpointRule().Update(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &EndpointRuleUpdateRes{ToEndpointRule(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
