package rest

import (
	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(router gin.IRoutes, service *sdk) {
	router.GET("/me", UseAccountGet(service))
}
