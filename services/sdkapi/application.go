package sdkapi

import (
	"github.com/gin-gonic/gin"
)

func RegisterApplicationRoutes(router gin.IRoutes, service *sdkapi) {
	router.POST("", UseApplicationCreate(service))
	router.PUT("/:app_id", UseApplicationUpdate(service))
	router.DELETE("/:app_id", UseApplicationDelete(service))

	router.GET("", UseApplicationList(service))
	router.GET("/:app_id", UseApplicationGet(service))
}
