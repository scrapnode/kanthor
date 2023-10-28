package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/gateway"
	ginmw "github.com/scrapnode/kanthor/gateway/gin/middlewares"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/sdk/usecase"
)

type MessagePutReq struct {
	Type string `json:"type" binding:"required" example:"testing.debug"`

	Body    map[string]interface{} `json:"body" binding:"required"`
	Headers map[string]string      `json:"headers"`
}

type MessagePutRes struct {
	Id string `json:"id"`
}

// UseMessagePut
// @Tags		message
// @Router		/application/{app_id}/message		[put]
// @Param		app_id								path		string			true	"application id"
// @Param		props								body		MessagePutReq	true	"message properties"
// @Success		201									{object}	MessagePutRes
// @Failure		default								{object}	gateway.Error
// @Security	BasicAuth
func UseMessagePut(logger logging.Logger, uc usecase.Sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req MessagePutReq
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

		appId := ginctx.Param("app_id")
		headers := entities.Header{}
		if len(req.Headers) > 0 {
			for k, v := range req.Headers {
				headers.Set(k, v)
			}
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		ws := ctx.Value(gateway.CtxWs).(*entities.Workspace)
		ucreq := &usecase.MessagePutReq{
			WsId:     ws.Id,
			Tier:     ws.Tier,
			AppId:    appId,
			Type:     req.Type,
			Body:     string(body),
			Headers:  headers,
			Metadata: entities.Metadata{entities.MetaMsgIdempotencyKey: ginctx.GetHeader(ginmw.HeaderIdempotencyKey)},
		}

		if err := ucreq.Validate(); err != nil {
			logger.Errorw(err.Error(), "data", utils.Stringify(ucreq))
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Message().Put(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &MessagePutRes{ucres.Message.Id}

		logger.Debugw("put message", "msg_id", ucres.Message.Id)
		ginctx.JSON(http.StatusCreated, res)
	}
}
