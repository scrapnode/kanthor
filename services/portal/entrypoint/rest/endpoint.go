package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterEndpointRoutes(router gin.IRoutes, service *portal) {
	router = router.
		Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc))).
		Use(middlewares.UseAuthz(service.infra.Authorizator))

	router.GET(":ep_id/message", UseEndpointListMessage(service))
	router.GET(":ep_id/message/:msg_id", UseEndpointGetMessage(service))
}
