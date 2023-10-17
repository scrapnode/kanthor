package portalapi

import (
	"github.com/gin-gonic/gin"
	ginmw "github.com/scrapnode/kanthor/infrastructure/gateway/gin/middlewares"
)

func RegisterWorkspaceRoutes(router gin.IRoutes, service *portalapi) {
	router = router.Use(ginmw.UseAuthz(service.infra.Authorizator))
	router.GET("/me", UseWorkspaceGet())
	router.PUT("/me", UseWorkspaceUpdate(service.logger, service.uc))
}
