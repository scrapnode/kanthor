package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type ApplicationGetMessageRes struct {
	*Message
} // @name ApplicationGetMessageRes

// UseApplicationGetMessage
// @Tags		application
// @Router		/application/{app_id}/message/{msg_id}		[get]
// @Success		200											{object}	ApplicationGetMessageRes
// @Failure		default										{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseApplicationGetMessage(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.ApplicationGetMessageIn{
			WsId:  ws.Id,
			AppId: ginctx.Param("app_id"),
			Id:    ginctx.Param("msg_id"),
		}
		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Application().GetMessage(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &ApplicationGetMessageRes{ToMessage(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
