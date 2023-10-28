package rest

import (
	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(router gin.IRoutes, service *sdk) {
	router.PUT("", UseMessagePut(service.logger, service.uc))
}
