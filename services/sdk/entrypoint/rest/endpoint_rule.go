package rest

import (
	"github.com/gin-gonic/gin"
)

func RegisterEndpointRuleRoutes(router gin.IRoutes, service *sdk) {
	router.POST("", UseEndpointRuleCreate(service))
	router.PUT("/:epr_id", UseEndpointRuleUpdate(service))
	router.DELETE("/:epr_id", UseEndpointRuleDelete(service))

	router.GET("", UseEndpointRuleList(service))
	router.GET("/:epr_id", UseEndpointRuleGet(service))
}
