package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type ApplicationCreateReq struct {
	Name string `json:"name" binding:"required"`
}

type ApplicationCreateRes struct {
	*entities.Application
}

// UseApplicationCreate Create new application in a workspace
//
//	@Summary		Create new application in a workspace
//	@Description	Create new application in a workspace
//	@Tags			application
//	@Router			/application	[post]
//	@Param			payload			body		ApplicationCreateReq	true	"application properties"
//	@Success		200				{object}	ApplicationCreateRes
//	@Failure		default			{object}	HttpError
//	@Security		BasicAuth
//	@in header
//	@name			Authorization
func UseApplicationCreate(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req ApplicationCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
			return
		}

		ctx := ginctx.MustGet("ctx").(context.Context)
		ucreq := &usecase.ApplicationCreateReq{Name: req.Name}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ucres, err := uc.Application().Create(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		res := &ApplicationCreateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
