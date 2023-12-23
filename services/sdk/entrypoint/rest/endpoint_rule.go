package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway/gin/middlewares"
)

func RegisterEndpointRuleRoutes(router gin.IRoutes, service *sdk) {
	router.
		Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc))).
		Use(middlewares.UseAuthz(service.infra.Authorizator))

	router.POST("", UseEndpointRuleCreate(service))
	router.PATCH("/:epr_id", UseEndpointRuleUpdate(service))
	router.DELETE("/:epr_id", UseEndpointRuleDelete(service))

	router.GET("", UseEndpointRuleList(service))
	router.GET("/:epr_id", UseEndpointRuleGet(service))
}
