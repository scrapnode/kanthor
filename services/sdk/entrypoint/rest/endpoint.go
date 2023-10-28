package rest

import (
	"github.com/gin-gonic/gin"
)

func RegisterEndpointRoutes(router gin.IRoutes, service *sdk) {
	router.POST("", UseEndpointCreate(service))
	router.PUT("/:ep_id", UseEndpointUpdate(service))
	router.DELETE("/:ep_id", UseEndpointDelete(service))

	router.GET("", UseEndpointList(service))
	router.GET("/:ep_id", UseEndpointGet(service))
}
