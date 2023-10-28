package rest

import (
	"github.com/gin-gonic/gin"
	ginmw "github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterWorkspaceRoutes(router gin.IRoutes, service *portal) {
	router = router.Use(ginmw.UseAuthz(service.infra.Authorizator))
	router.GET("/me", UseWorkspaceGet())
	router.PUT("/me", UseWorkspaceUpdate(service))
}
