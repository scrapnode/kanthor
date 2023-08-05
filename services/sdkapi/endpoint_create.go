package sdkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type EndpointCreateReq struct {
	Name string `json:"name" binding:"required"`

	SecretKey string `json:"secret_key" binding:"omitempty,min=16,max=32"`
	Method    string `json:"method" binding:"required,oneof=POST PUT"`
	Uri       string `json:"uri" binding:"required,uri" example:"https://example.com"`
}

type EndpointCreateRes struct {
	*entities.Endpoint
}

// UseEndpointCreate
// @Tags		endpoint
// @Router		/application/{app_id}/endpoint		[post]
// @Param		app_id								path		string				true	"application id"
// @Param		props								body		EndpointCreateReq	true	"endpoint properties"
// @Success		201									{object}	EndpointCreateRes
// @Failure		default								{object}	gateway.Error
// @Security	BasicAuth
// @in header
// @name		Authorization
func UseEndpointCreate(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet("ctx").(context.Context)
		appId := ginctx.Param("app_id")
		ucreq := &usecase.EndpointCreateReq{
			AppId:     appId,
			Name:      req.Name,
			SecretKey: req.SecretKey,
			Method:    req.Method,
			Uri:       req.Uri,
		}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Endpoint().Create(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointCreateRes{ucres.Doc}
		ginctx.JSON(http.StatusCreated, res)
	}
}
