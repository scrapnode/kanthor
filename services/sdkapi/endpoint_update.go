package sdkapi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/pkg/utils"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

type EndpointUpdateReq struct {
	Name string `json:"name" binding:"required"`
}

type EndpointUpdateRes struct {
	*entities.Endpoint
}

// UseEndpointUpdate
// @Tags		endpoint
// @Router		/application/{app_id}/endpoint/{ep_id}	[put]
// @Param		app_id									path		string					true	"application id"
// @Param		ep_id									path		string					true	"endpoint id"
// @Param		props									body		EndpointUpdateReq		true	"endpoint properties"
// @Success		200										{object}	EndpointUpdateRes
// @Failure		default									{object}	gateway.Error
// @Security	BasicAuth
func UseEndpointUpdate(service *sdkapi) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		appId := ginctx.Param("app_id")
		id := ginctx.Param("ep_id")
		ucreq := &usecase.EndpointUpdateReq{AppId: appId, Id: id, Name: req.Name}
		if err := ucreq.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := service.uc.Endpoint().Update(ctx, ucreq)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointUpdateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
