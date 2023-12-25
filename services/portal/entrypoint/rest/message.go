package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterMessageRoutes(router gin.IRoutes, service *portal) {
	router = router.
		Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc))).
		Use(middlewares.UseAuthz(service.infra.Authorizator))

	router.GET("", UseMessageList(service))
	router.GET("/:msg_id", UseMessageGet(service))
}
