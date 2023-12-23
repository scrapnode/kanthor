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
	*EndpointRule
}

// UseEndpointRuleUpdate
// @Tags		endpoint rule
// @Router		/rule/{epr_id}	[patch]
// @Param		epr_id			path		string					true	"rule id"
// @Param		props			body		EndpointRuleUpdateReq	true	"rule properties"
// @Success		200				{object}	EndpointRuleUpdateRes
// @Failure		default			{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointRuleUpdate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointRuleUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointRuleUpdateIn{
			WsId: ws.Id,
			Id:   ginctx.Param("epr_id"),
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

		res := &EndpointRuleUpdateRes{ToEndpointRule(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
