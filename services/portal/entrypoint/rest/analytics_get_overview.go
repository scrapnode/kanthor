package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type AnalyticsGetOverviewRes struct {
	CredentialsCount int64 `json:"credentials_count"`
	ApplicationCount int64 `json:"application_count"`
	EndpointCount    int64 `json:"endpoint_count"`
} // @name AnalyticsGetOverviewRes

// UseAnalyticsGetOverview
// @Tags		analytics
// @Router		/analytics/overview		[get]
// @Success		200						{object}	AnalyticsGetOverviewRes
// @Failure		default					{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseAnalyticsGetOverview(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.AnalyticsGetOverviewIn{
			WsId: ws.Id,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Analytics().GetOverview(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &AnalyticsGetOverviewRes{
			CredentialsCount: out.CredentialsCount,
			ApplicationCount: out.ApplicationCount,
			EndpointCount:    out.CredentialsCount,
		}
		ginctx.JSON(http.StatusOK, res)
	}
}
