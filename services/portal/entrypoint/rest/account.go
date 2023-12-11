package rest

import (
	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(router gin.IRoutes, service *portal) {
	router.POST("/setup", UseAccountSetup(service))
	router.GET("/me", UseAccountGet(service))
}
