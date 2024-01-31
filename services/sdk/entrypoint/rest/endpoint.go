package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterEndpointRoutes(router gin.IRoutes, service *sdk) {
	router.Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc)))

	router.POST("", UseEndpointCreate(service))
	router.PATCH("/:ep_id", UseEndpointUpdate(service))
	router.DELETE("/:ep_id", UseEndpointDelete(service))

	router.GET("", UseEndpointList(service))
	router.GET("/:ep_id", UseEndpointGet(service))
}
