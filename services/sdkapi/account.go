package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(router gin.IRoutes, service *sdkapi) {
	router.GET("/me", UseAccountGet(service))
}
