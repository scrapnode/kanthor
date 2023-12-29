package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterRequestRoutes(router gin.IRoutes, service *portal) {
	router = router.
		Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc))).
		Use(middlewares.UseAuthz(service.infra.Authorizator))

	router.GET("", UseRequestList(service))
	router.GET("/:req_id", UseRequestGet(service))
}
