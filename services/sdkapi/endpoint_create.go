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

type endpointCreateReq struct {
	Name string `json:"name" binding:"required"`

	SecretKey string `json:"secret_key" binding:"omitempty,min=16,max=32"`
	Method    string `json:"method" binding:"required,oneof=POST PUT"`
	Uri       string `json:"uri" binding:"required,uri"`
}

type endpointCreateRes struct {
	*entities.Endpoint
}

func UseEndpointCreate(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req endpointCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
			return
		}

		ctx := ginctx.MustGet("ctx").(context.Context)
		app := ginctx.MustGet("app").(*entities.Application)

		ucreq := &usecase.EndpointCreateReq{
			AppId:     app.Id,
			Name:      req.Name,
			SecretKey: req.SecretKey,
			Method:    req.Method,
			Uri:       req.Uri,
		}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ucres, err := uc.Endpoint().Create(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "oops, something went wrong"})
			return
		}

		res := &endpointCreateRes{ucres.Doc}
		ginctx.JSON(http.StatusOK, res)
	}
}
