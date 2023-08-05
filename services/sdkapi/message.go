package sdkapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

func UseMessageRoutes(router *gin.RouterGroup, logger logging.Logger, validator validator.Validator, uc usecase.Sdk) {
	router.PUT("", UseMessagePut(logger, validator, uc))
}
