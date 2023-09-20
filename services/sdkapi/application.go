package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func RegisterApplicationRoutes(router gin.IRoutes, service *sdkapi) {
	router.POST("", UseApplicationCreate(service.logger, service.uc))
	router.PUT("/:app_id", UseApplicationUpdate(service.logger, service.uc))
	router.DELETE("/:app_id", UseApplicationDelete(service.logger, service.uc))

	router.GET("", UseApplicationList(service.logger, service.uc))
	router.GET("/:app_id", UseApplicationGet(service.logger, service.uc))
}
