package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointCreateReq struct {
	AppId string `json:"app_id"`
	Name  string `json:"name" default:"POST httpbin"`

	Method string `json:"method" example:"POST"`
	Uri    string `json:"uri" example:"https://httpbin.org/post"`
} // @name EndpointCreateReq

type EndpointCreateRes struct {
	*Endpoint
	// To make the UI become friendly we will return the secret key after user create the new endpoint
	// but we don't want to return that key everytime user request for the endpoint
	// user must have specific permission to reveal the secret key of an endpoint
	SecretKey string `json:"secret_key"`
} // @name EndpointCreateRes

// UseEndpointCreate
// @Tags			endpoint
// @Router		/endpoint			[post]
// @Param			payload				body			EndpointCreateReq	true	"endpoint payload"
// @Success		201						{object}	EndpointCreateRes
// @Failure		default				{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointCreate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req EndpointCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointCreateIn{
			WsId:      ws.Id,
			AppId:     req.AppId,
			Name:      req.Name,
			SecretKey: utils.RandomString(32),
			Method:    req.Method,
			Uri:       req.Uri,
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Endpoint().Create(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &EndpointCreateRes{ToEndpoint(out.Doc), out.Doc.SecretKey}
		ginctx.JSON(http.StatusCreated, res)
	}
}
