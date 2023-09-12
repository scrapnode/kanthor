package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func UseAccountRoutes(router gin.IRoutes, service *sdkapi) {
	router.GET("/me", UseAccountGet(service))
}
