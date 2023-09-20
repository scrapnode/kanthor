package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func RegisterEndpointRuleRoutes(router gin.IRoutes, service *sdkapi) {
	router.POST("", UseEndpointRuleCreate(service.logger, service.uc))
	router.PUT("/:epr_id", UseEndpointRuleUpdate(service.logger, service.uc))
	router.DELETE("/:epr_id", UseEndpointRuleDelete(service.logger, service.uc))

	router.GET("", UseEndpointRuleList(service.logger, service.uc))
	router.GET("/:epr_id", UseEndpointRuleGet(service.logger, service.uc))
}
