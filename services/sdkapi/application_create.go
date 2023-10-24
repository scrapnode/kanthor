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
func UseApplicationCreate(service *sdkapi) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req ApplicationCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		ucreq := &usecase.ApplicationCreateReq{Name: req.Name}
		if err := ucreq.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := service.uc.Application().Create(ctx, ucreq)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &ApplicationCreateRes{ucres.Doc}
		ginctx.JSON(http.StatusCreated, res)
	}
}
