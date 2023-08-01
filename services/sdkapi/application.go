package sdkapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

func UseApplication(router *gin.RouterGroup, logger logging.Logger, uc usecase.Sdk) {
	router.POST("", UseApplicationCreate(logger, uc))
}
