package portalapi

import (
	"github.com/gin-gonic/gin"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
)

func RegisterWorkspaceCredentialsRoutes(router gin.IRoutes, service *portalapi) {
	router = router.Use(ginmw.UseAuthz(service.infra.Authorizator))
	router.GET("", UseWorkspaceCredentialsList(service.logger, service.uc))
	router.POST("", UseWorkspaceCredentialsCreate(service.logger, service.uc))
	router.GET("/:wsc_id", UseWorkspaceCredentialsGet(service.logger, service.uc))
	router.PUT("/:wsc_id", UseWorkspaceCredentialsUpdate(service.logger, service.uc))
	router.PUT("/:wsc_id/expiration", UseWorkspaceCredentialsExpire(service.logger, service.uc))
}
