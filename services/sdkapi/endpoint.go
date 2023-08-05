package sdkapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

func UseEndpointRoutes(router *gin.RouterGroup, logger logging.Logger, validator validator.Validator, uc usecase.Sdk) {
	router.POST("", UseEndpointCreate(logger, validator, uc))
	router.PATCH("/:ep_id", UseEndpointUpdate(logger, validator, uc))
	router.DELETE("/:ep_id", UseEndpointDelete(logger, validator, uc))

	router.GET("", UseEndpointList(logger, validator, uc))
	router.GET("/:ep_id", UseEndpointGet(logger, validator, uc))
}
