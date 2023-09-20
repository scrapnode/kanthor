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

type ApplicationUpdateReq struct {
	Name string `json:"name" binding:"required"`
}

type ApplicationUpdateRes struct {
	*entities.Application
}

// UseApplicationUpdate
// @Tags		application
// @Router		/application/{app_id}	[put]
// @Param		app_id					path		string					true	"application id"
// @Param		props					body		ApplicationUpdateReq	true	"application properties"
// @Success		200						{object}	ApplicationUpdateRes
// @Failure		default					{object}	gateway.Error
// @Security	BasicAuth
func UseApplicationUpdate(logger logging.Logger, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req ApplicationUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		id := ginctx.Param("app_id")
		ucreq := &usecase.ApplicationUpdateReq{Id: id, Name: req.Name}
		if err := ucreq.Validate(); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Application().Update(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &ApplicationUpdateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
