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

type EndpointRuleGetRes struct {
	*EndpointRule
}

// UseEndpointRuleGet
// @Tags		endpoint rule
// @Router		/rule/{epr_id}	[get]
// @Param		epr_id			path		string					true	"rule id"
// @Param		ep_id			query		string					true	"endpoint id"
// @Success		200				{object}	EndpointRuleGetRes
// @Failure		default			{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointRuleGet(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointRuleGetIn{
			WsId: ws.Id,
			Id:   ginctx.Param("epr_id"),
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

		res := &EndpointRuleGetRes{ToEndpointRule(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
