package sdkapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
)

func UseEndpointRuleRoutes(router *gin.RouterGroup, logger logging.Logger, validator validator.Validator, uc usecase.Sdk) {
	router.POST("", UseEndpointRuleCreate(logger, validator, uc))
	router.PATCH("/:epr_id", UseEndpointRuleUpdate(logger, validator, uc))
	router.DELETE("/:epr_id", UseEndpointRuleDelete(logger, validator, uc))

	router.GET("", UseEndpointRuleList(logger, validator, uc))
	router.GET("/:epr_id", UseEndpointRuleGet(logger, validator, uc))
}
