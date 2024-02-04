package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type EndpointGetSecretRes struct {
	SecretKey string `json:"secret_key"`
} // @name EndpointGetSecretRes

// UseEndpointGetSecret
// @Tags			endpoint
// @Router		/endpoint/{ep_id}/secret	[get]
// @Param			ep_id											path		string					true	"endpoint id"
// @Success		200												{object}	EndpointGetSecretRes
// @Failure		default										{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseEndpointGetSecret(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.EndpointGetIn{
			WsId: ws.Id,
			Id:   ginctx.Param("ep_id"),
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Endpoint().Get(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &EndpointGetSecretRes{SecretKey: out.Doc.SecretKey}
		ginctx.JSON(http.StatusOK, res)
	}
}
