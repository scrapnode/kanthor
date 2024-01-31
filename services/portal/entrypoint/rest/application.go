package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterApplicationRoutes(router gin.IRoutes, service *portal) {
	router = router.Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc)))

	router.GET(":app_id/message", UseApplicationListMessage(service))
	router.GET(":app_id/message/:msg_id", UseApplicationGetMessage(service))
}
