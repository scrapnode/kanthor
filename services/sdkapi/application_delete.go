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

type ApplicationDeleteRes struct {
	*entities.Application
}

// UseApplicationDelete
// @Tags		application
// @Router		/application/{app_id}	[delete]
// @Param		app_id					path		string					true	"application id"
// @Success		200						{object}	ApplicationDeleteRes
// @Failure		default					{object}	gateway.Error
// @Security	BasicAuth
func UseApplicationDelete(logger logging.Logger, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		id := ginctx.Param("app_id")
		ucreq := &usecase.ApplicationDeleteReq{Id: id}
		if err := ucreq.Validate(); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Application().Delete(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &ApplicationDeleteRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
