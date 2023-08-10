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

type ApplicationCreateReq struct {
	Name string `json:"name" binding:"required"`
}

type ApplicationCreateRes struct {
	*entities.Application
}

// UseApplicationCreate
// @Tags		application
// @Router		/application		[post]
// @Param		props				body		ApplicationCreateReq	true	"application properties"
// @Success		201					{object}	ApplicationCreateRes
// @Failure		default				{object}	gateway.Error
// @Security	BasicAuth
func UseApplicationCreate(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req ApplicationCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet("ctx").(context.Context)
		ucreq := &usecase.ApplicationCreateReq{Name: req.Name}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Application().Create(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &ApplicationCreateRes{ucres.Doc}
		ginctx.JSON(http.StatusCreated, res)
	}
}
