package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointRuleCreateReq struct {
	EpId string `json:"ep_id"`
	Name string `json:"name" `

	Priority            int32  `json:"priority"`
	Exclusionary        bool   `json:"exclusionary"`
	ConditionSource     string `json:"condition_source"`
	ConditionExpression string `json:"condition_expression"`
} // @name EndpointRuleCreateReq

type EndpointRuleCreateRes struct {
	*EndpointRule
} // @name EndpointRuleCreateRes

// UseEndpointRuleCreate
// @Tags		endpoint rule
// @Router		/rule		[post]
// @Param		payload		body		EndpointRuleCreateReq	true	"rule payload"
// @Success		201			{object}	EndpointRuleCreateRes
// @Failure		default		{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointRuleCreate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointRuleCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointRuleCreateIn{
			WsId:                ws.Id,
			EpId:                req.EpId,
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

		out, err := service.uc.EndpointRule().Create(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &EndpointRuleCreateRes{ToEndpointRule(out.Doc)}
		ginctx.JSON(http.StatusCreated, res)
	}
}
