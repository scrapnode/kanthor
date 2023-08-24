package portalapi

import (
	"github.com/gin-gonic/gin"
)

func UseAccountRoutes(router *gin.RouterGroup) {
	router.GET("/me", UseAccountGet())
}
