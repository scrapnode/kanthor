package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterMessageRoutes(router gin.IRoutes, service *sdk) {
	router.Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc)))

	router.POST("", UseMessageCreate(service))
}
