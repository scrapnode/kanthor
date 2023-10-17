package sdkapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

type EndpointGetRes struct {
	*entities.Endpoint
}

// UseEndpointGet
// @Tags		endpoint
// @Router		/application/{app_id}/endpoint/{ep_id}	[get]
// @Param		app_id									path		string					true	"application id"
// @Param		ep_id									path		string					true	"endpoint id"
// @Success		200										{object}	EndpointGetRes
// @Failure		default									{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointGet(logger logging.Logger, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		appId := ginctx.Param("app_id")
		id := ginctx.Param("ep_id")
		ucreq := &usecase.EndpointGetReq{AppId: appId, Id: id}
		if err := ucreq.Validate(); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Endpoint().Get(ctx, ucreq)
		if err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointGetRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
