package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func UseEndpointRoutes(router gin.IRoutes, service *sdkapi) {
	router.POST("", UseEndpointCreate(service.logger, service.validator, service.uc))
	router.PUT("/:ep_id", UseEndpointUpdate(service.logger, service.validator, service.uc))
	router.DELETE("/:ep_id", UseEndpointDelete(service.logger, service.validator, service.uc))

	router.GET("", UseEndpointList(service.logger, service.validator, service.uc))
	router.GET("/:ep_id", UseEndpointGet(service.logger, service.validator, service.uc))
}
