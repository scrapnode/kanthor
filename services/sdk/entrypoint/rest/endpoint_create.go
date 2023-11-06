package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
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
func UseEndpointCreate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		appId := ginctx.Param("app_id")
		in := &usecase.EndpointCreateIn{
			AppId:     appId,
			Name:      req.Name,
			SecretKey: req.SecretKey,
			Method:    req.Method,
			Uri:       req.Uri,
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Endpoint().Create(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &EndpointCreateRes{out.Doc}
		ginctx.JSON(http.StatusCreated, res)
	}
}
