package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type endpointDeleteRes struct {
	*entities.Endpoint
}

func UseEndpointDelete(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet("ctx").(context.Context)
		appId := ginctx.Param("app_id")
		id := ginctx.Param("ep_id")
		ucreq := &usecase.EndpointDeleteReq{AppId: appId, Id: id}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ucres, err := uc.Endpoint().Delete(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		res := &endpointDeleteRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
