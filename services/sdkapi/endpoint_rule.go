package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func UseEndpointRuleRoutes(router gin.IRoutes, service *sdkapi) {
	router.POST("", UseEndpointRuleCreate(service.logger, service.validator, service.uc))
	router.PUT("/:epr_id", UseEndpointRuleUpdate(service.logger, service.validator, service.uc))
	router.DELETE("/:epr_id", UseEndpointRuleDelete(service.logger, service.validator, service.uc))

	router.GET("", UseEndpointRuleList(service.logger, service.validator, service.uc))
	router.GET("/:epr_id", UseEndpointRuleGet(service.logger, service.validator, service.uc))
}
