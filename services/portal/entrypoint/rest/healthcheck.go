package rest

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterHealthcheck(router gin.IRoutes, service *portal) {
	router.GET("/", func(ginctx *gin.Context) {
		host, _ := os.Hostname()
		ginctx.JSON(http.StatusOK, gin.H{"host": host, "service": "portal"})
	})

	router.GET("/readiness", func(ginctx *gin.Context) {
		if err := service.infra.Readiness(); err != nil {
			ginctx.Status(http.StatusInternalServerError)
			return
		}

		if err := service.db.Readiness(); err != nil {
			ginctx.Status(http.StatusInternalServerError)
			return
		}

		ginctx.String(http.StatusOK, "ready")
	})
	router.GET("/liveness", func(ginctx *gin.Context) {
		if err := service.infra.Liveness(); err != nil {
			ginctx.Status(http.StatusInternalServerError)
			return
		}

		if err := service.db.Liveness(); err != nil {
			ginctx.Status(http.StatusInternalServerError)
			return
		}

		ginctx.String(http.StatusOK, "live")
	})
}
