package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterApplicationRoutes(router gin.IRoutes, service *sdk) {
	router.Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc)))

	router.POST("", UseApplicationCreate(service))
	router.PATCH("/:app_id", UseApplicationUpdate(service))
	router.DELETE("/:app_id", UseApplicationDelete(service))

	router.GET("", UseApplicationList(service))
	router.GET("/:app_id", UseApplicationGet(service))
}
