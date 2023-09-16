package portalapi

import (
	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(router gin.IRoutes, service *portalapi) {
	router.PUT("/me", UseAccountSetup(service.logger, service.validator, service.uc))
	router.GET("/me", UseAccountGet())
}
