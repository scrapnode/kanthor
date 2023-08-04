package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

func UseApplicationRoutes(router *gin.RouterGroup, logger logging.Logger, validator validator.Validator, uc usecase.Sdk) {
	router.POST("", UseApplicationCreate(logger, validator, uc))
	router.PATCH("/:app_id", UseApplicationUpdate(logger, validator, uc))
	router.DELETE("/:app_id", UseApplicationDelete(logger, validator, uc))

	router.GET("", UseApplicationList(logger, validator, uc))
	router.GET("/:app_id", UseApplication(logger, validator, uc), UseApplicationGet())

	UseEndpointRoutes(router.Group("/:app_id/endpoint"), logger, validator, uc)
}

func UseApplication(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet("ctx").(context.Context)
		id := ginctx.Param("app_id")
		ucreq := &usecase.ApplicationGetReq{Id: id}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ucres, err := uc.Application().Get(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		ginctx.Set("app", ucres.Doc)
		ginctx.Next()
	}
}
