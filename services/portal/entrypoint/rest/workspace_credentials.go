package rest

import (
	"github.com/gin-gonic/gin"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
)

func RegisterWorkspaceCredentialsRoutes(router gin.IRoutes, service *portal) {
	router = router.Use(ginmw.UseAuthz(service.infra.Authorizator))
	router.GET("", UseWorkspaceCredentialsList(service))
	router.POST("", UseWorkspaceCredentialsCreate(service))
	router.GET("/:wsc_id", UseWorkspaceCredentialsGet(service))
	router.PUT("/:wsc_id", UseWorkspaceCredentialsUpdate(service))
	router.PUT("/:wsc_id/expiration", UseWorkspaceCredentialsExpire(service))
}
