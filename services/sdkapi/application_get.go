package sdkapi

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"net/http"
)

type applicationGetRes struct {
	*entities.Application
}

func UseApplicationGet() gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		app, ok := ginctx.Get("app")
		if !ok {
			ginctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "application is not found"})
			return
		}

		res := &applicationGetRes{app.(*entities.Application)}
		ginctx.JSON(http.StatusOK, res)
	}
}
