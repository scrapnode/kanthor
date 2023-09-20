package portalapi

import (
	"github.com/gin-gonic/gin"
)

func RegisterWorkspaceRoutes(router gin.IRoutes, service *portalapi) {
	router.GET("/me", UseWorkspaceGet())
	router.PUT("/me", UseWorkspaceUpdate(service.logger, service.uc))
}
