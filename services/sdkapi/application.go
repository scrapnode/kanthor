package sdkapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

func UseApplication(router *gin.RouterGroup, logger logging.Logger, uc usecase.Sdk) {
	router.POST("", UseApplicationCreate(logger, uc))
	router.PATCH("/:app_id", UseApplicationUpdate(logger, uc))
	router.DELETE("/:app_id", UseApplicationDelete(logger, uc))

	router.GET("", UseApplicationList(logger, uc))
	router.GET("/:app_id", UseApplicationGet(logger, uc))
}
