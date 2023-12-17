package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterWorkspaceRoutes(router gin.IRoutes, service *portal) {
	// exception: do not required selected workspace
	router.GET("", UseWorkspaceList(service))

	router = router.
		Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc))).
		Use(middlewares.UseAuthz(service.infra.Authorizator))

	router.GET("/me", UseWorkspaceGet(service))
	router.PATCH("/me", UseWorkspaceUpdate(service))

	router.GET("/me/credentials", UseWorkspaceCredentialsList(service))
	router.POST("/me/credentials", UseWorkspaceCredentialsCreate(service))
	router.GET("/me/credentials/:wsc_id", UseWorkspaceCredentialsGet(service))
	router.PATCH("/me/credentials/:wsc_id", UseWorkspaceCredentialsUpdate(service))
	router.PUT("/me/credentials/:wsc_id/expiration", UseWorkspaceCredentialsExpire(service))
}
