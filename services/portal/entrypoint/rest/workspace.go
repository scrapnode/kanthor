package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterWorkspaceRoutes(router gin.IRoutes, service *portal) {
	router = router.
		Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc))).
		Use(middlewares.UseAuthz(service.infra.Authorizator))

	router.GET("", UseWorkspaceGet())
	router.PUT("", UseWorkspaceUpdate(service))
}
