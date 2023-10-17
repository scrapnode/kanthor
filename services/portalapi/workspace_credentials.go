package portalapi

import (
	"github.com/gin-gonic/gin"
)

func RegisterWorkspaceCredentialsRoutes(router gin.IRoutes, service *portalapi) {
	router.GET("", UseWorkspaceCredentialsList(service.logger, service.uc))
	router.POST("", UseWorkspaceCredentialsCreate(service.logger, service.uc, service.infra.Coordinator))
	router.GET("/:wsc_id", UseWorkspaceCredentialsGet(service.logger, service.uc))
	router.PUT("/:wsc_id", UseWorkspaceCredentialsUpdate(service.logger, service.uc))
	router.PUT("/:wsc_id/expiration", UseWorkspaceCredentialsExpire(service.logger, service.uc, service.infra.Coordinator))
}
