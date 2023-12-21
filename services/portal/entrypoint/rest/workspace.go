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
}
