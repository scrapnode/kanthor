package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointDeleteRes struct {
	*entities.Endpoint
}

// UseEndpointDelete
// @Tags		endpoint
// @Router		/application/{app_id}/endpoint/{ep_id}	[delete]
// @Param		app_id									path		string					true	"application id"
// @Param		ep_id									path		string					true	"endpoint id"
// @Success		200										{object}	EndpointDeleteRes
// @Failure		default									{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointDelete(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		appId := ginctx.Param("app_id")
		id := ginctx.Param("ep_id")
		ucreq := &usecase.EndpointDeleteReq{AppId: appId, Id: id}
		if err := ucreq.Validate(); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := service.uc.Endpoint().Delete(ctx, ucreq)
		if err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointDeleteRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
