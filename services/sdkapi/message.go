package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func UseMessageRoutes(router *gin.RouterGroup, service *sdkapi) {
	router.PUT("", UseMessagePut(service.logger, service.validator, service.uc))
}
