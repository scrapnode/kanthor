package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type ApplicationDeleteReq struct {
	Name string `json:"name" binding:"required"`
}

type ApplicationDeleteRes struct {
	*entities.Application
}

func UseApplicationDelete(logger logging.Logger, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req ApplicationDeleteReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
			return
		}

		ctx := ginctx.MustGet("ctx").(context.Context)
		ucreq := &usecase.ApplicationDeleteReq{Id: ginctx.Param("app_id")}
		ucres, err := uc.Application().Delete(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		res := &ApplicationDeleteRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
