package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type ApplicationGetRes struct {
	*entities.Application
}

// UseApplicationGet
// @Tags		application
// @Router		/application/{app_id}	[get]
// @Param		app_id					path		string					true	"application id"
// @Success		200						{object}	ApplicationGetRes
// @Failure		default					{object}	gateway.Error
// @Security	BasicAuth
func UseApplicationGet(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet("ctx").(context.Context)
		id := ginctx.Param("app_id")
		ucreq := &usecase.ApplicationGetReq{Id: id}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Application().Get(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &ApplicationGetRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
