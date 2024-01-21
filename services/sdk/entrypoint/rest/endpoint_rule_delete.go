package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointRuleDeleteRes struct {
	*EndpointRule
} // @name EndpointRuleDeleteRes

// UseEndpointRuleDelete
// @Tags		endpoint rule
// @Router		/rule/{epr_id}	[delete]
// @Param		epr_id			path		string					true	"rule id"
// @Success		200				{object}	EndpointRuleDeleteRes
// @Failure		default			{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointRuleDelete(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointRuleDeleteIn{
			WsId: ws.Id,
			Id:   ginctx.Param("epr_id"),
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.EndpointRule().Delete(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &EndpointRuleDeleteRes{ToEndpointRule(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
