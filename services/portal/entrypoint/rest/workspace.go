package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

type WorkspaceCredentials struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`

	WsId      string `json:"ws_id"`
	Name      string `json:"name"`
	ExpiredAt int64  `json:"expired_at"`
}

func RegisterWorkspaceRoutes(router gin.IRoutes, service *portal) {
	// exception: do not required selected workspace
	router.GET("", UseWorkspaceList(service))

	router = router.
		Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc))).
		Use(middlewares.UseAuthz(service.infra.Authorizator))

	router.GET("/me", UseWorkspaceGet(service))
	router.PUT("/me", UseWorkspaceUpdate(service))

	router.GET("/me/credentials", UseWorkspaceCredentialsList(service))
	router.POST("/me/credentials", UseWorkspaceCredentialsCreate(service))
	router.GET("/me/credentials/:wsc_id", UseWorkspaceCredentialsGet(service))
	router.PATCH("/me/credentials/:wsc_id", UseWorkspaceCredentialsUpdate(service))
	router.PUT("/me/credentials/:wsc_id/expiration", UseWorkspaceCredentialsExpire(service))
}
