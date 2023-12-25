package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

type MessageGetRes struct {
	*Message
} // @name MessageGetRes

// UseMessageGet
// @Tags		message
// @Router		/message/{msg_id}	[get]
// @Param		app_id					query		string			true	"application id"
// @Param		msg_id					path		string			true	"message id"
// @Success		200						{object}	MessageGetRes
// @Failure		default					{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseMessageGet(service *portal) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)

		in := &usecase.MessageGetIn{
			WsId:  ws.Id,
			AppId: ginctx.Query("app_id"),
			Id:    ginctx.Param("msg_id"),
		}
		if err := in.Validate(); err != nil {
			service.logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := service.uc.Message().Get(ctx, in)
		if err != nil {
			service.logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &MessageGetRes{ToMessage(out.Doc)}
		ginctx.JSON(http.StatusOK, res)
	}
}
