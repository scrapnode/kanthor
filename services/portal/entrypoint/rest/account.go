package rest

import (
	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(router gin.IRoutes, service *portal) {
	router.GET("", UseAccountGet(service))
}
