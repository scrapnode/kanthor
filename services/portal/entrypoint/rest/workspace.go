package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterWorkspaceRoutes(router gin.IRoutes, service *portal) {
	router = router.
		Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc))).
		Use(middlewares.UseAuthz(service.infra.Authorizator))

	router.GET("/me", UseWorkspaceGet())
	router.PUT("/me", UseWorkspaceUpdate(service))

	router.GET("/em/credentials", UseWorkspaceCredentialsList(service))
	router.POST("/em/credentials", UseWorkspaceCredentialsCreate(service))
	router.GET("/em/credentials/:wsc_id", UseWorkspaceCredentialsGet(service))
	router.PUT("/em/credentials/:wsc_id", UseWorkspaceCredentialsUpdate(service))
	router.PUT("/em/credentials/:wsc_id/expiration", UseWorkspaceCredentialsExpire(service))
}
