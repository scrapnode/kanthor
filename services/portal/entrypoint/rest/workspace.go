package rest

import (
	"github.com/gin-gonic/gin"
)

func RegisterWorkspaceRoutes(router gin.IRoutes, service *portal) {
	router.GET("", UseWorkspaceList(service))
	router.POST("", UseWorkspaceCreate(service))
	router.GET("/:ws_id", UseWorkspaceGet(service))
	router.PATCH("/:ws_id", UseWorkspaceUpdate(service))
	router.GET("/:ws_id/transfer", UseWorkspaceExport(service))
	router.POST("/:ws_id/transfer", UseWorkspaceImport(service))
}
