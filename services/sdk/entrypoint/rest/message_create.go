package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	ginmw "github.com/scrapnode/kanthor/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type MessageCreateReq struct {
	AppId string `json:"app_id"`
	Type  string `json:"type" example:"testing.debug"`

	Body    map[string]interface{} `json:"body"`
	Headers map[string]string      `json:"headers"`
} // @name MessageCreateReq

type MessageCreateRes struct {
	Id string `json:"id"`
} // @name MessageCreateRes

// UseMessageCreate
// @Tags		message
// @Router		/message		[post]
// @Param		payload			body		MessageCreateReq	true	"message payload"
// @Success		201				{object}	MessageCreateRes
// @Failure		default			{object}	gateway.Error
// @Security	Authorization
// @Security	WorkspaceId
func UseMessageCreate(logger logging.Logger, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req MessageCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed request"))
			return
		}

		body, err := json.Marshal(req.Body)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("malformed body"))
			return
		}

		headers := entities.Header{}
		if len(req.Headers) > 0 {
			for k, v := range req.Headers {
				headers.Set(k, v)
			}
		}

		ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ws := ctx.Value(gateway.CtxWorkspace).(*entities.Workspace)
		in := &usecase.MessageCreateIn{
			WsId:     ws.Id,
			Tier:     ws.Tier,
			AppId:    req.AppId,
			Type:     req.Type,
			Body:     string(body),
			Headers:  headers,
			Metadata: entities.Metadata{entities.MetaMsgIdempotencyKey: ginctx.GetHeader(ginmw.HeaderIdempotencyKey)},
		}

		if err := in.Validate(); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(in))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		out, err := uc.Message().Create(ctx, in)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &MessageCreateRes{out.Message.Id}

		logger.Debugw("put message", "msg_id", out.Message.Id)
		ginctx.JSON(http.StatusCreated, res)
	}
}
