package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterWorkspaceCredentialsRoutes(router gin.IRoutes, service *portal) {
	router = router.Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc)))

	router.GET("", UseWorkspaceCredentialsList(service))
	router.POST("", UseWorkspaceCredentialsCreate(service))
	router.GET("/:wsc_id", UseWorkspaceCredentialsGet(service))
	router.PATCH("/:wsc_id", UseWorkspaceCredentialsUpdate(service))
	router.PUT("/:wsc_id/expiration", UseWorkspaceCredentialsExpire(service))
}
