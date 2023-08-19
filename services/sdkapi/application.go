package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func UseApplicationRoutes(router *gin.RouterGroup, service *sdkapi) {
	router.POST("", UseApplicationCreate(service.logger, service.validator, service.uc))
	router.PUT("/:app_id", UseApplicationUpdate(service.logger, service.validator, service.uc))
	router.DELETE("/:app_id", UseApplicationDelete(service.logger, service.validator, service.uc))

	router.GET("", UseApplicationList(service.logger, service.validator, service.uc))
	router.GET("/:app_id", UseApplicationGet(service.logger, service.validator, service.uc))
}
