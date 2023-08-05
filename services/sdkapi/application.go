package sdkapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

func UseApplicationRoutes(router *gin.RouterGroup, logger logging.Logger, validator validator.Validator, uc usecase.Sdk) {
	router.POST("", UseApplicationCreate(logger, validator, uc))
	router.PATCH("/:app_id", UseApplicationUpdate(logger, validator, uc))
	router.DELETE("/:app_id", UseApplicationDelete(logger, validator, uc))

	router.GET("", UseApplicationList(logger, validator, uc))
	router.GET("/:app_id", UseApplicationGet(logger, validator, uc))
}
