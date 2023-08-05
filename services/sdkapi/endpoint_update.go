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

type EndpointUpdateReq struct {
	Name string `json:"name" binding:"required"`
}

type endpointUpdateRes struct {
	*entities.Endpoint
}

func UseEndpointUpdate(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointUpdateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ctx := ginctx.MustGet("ctx").(context.Context)
		app := ginctx.MustGet("app").(*entities.Application)

		id := ginctx.Param("ep_id")
		ucreq := &usecase.EndpointUpdateReq{AppId: app.Id, Id: id, Name: req.Name}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Endpoint().Update(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &endpointUpdateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
