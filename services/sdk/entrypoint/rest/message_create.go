package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/internal/entities"
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
// @Failure		default			{object}	gateway.Err
// @Security	Authorization
// @Security	WorkspaceId
func UseMessageCreate(service *sdk) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		var req MessageCreateReq
		if err := ginctx.ShouldBindJSON(&req); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		body, err := json.Marshal(req.Body)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
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
			Metadata: entities.Metadata{},
		}
		if ct := in.Headers.Get("Content-Type"); ct == "" {
			in.Headers.Set("Content-Type", ginctx.Request.Header.Get("Content-Type"))
		}

		if err := in.Validate(); err != nil {
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.Error(err))
			return
		}

		out, err := service.uc.Message().Create(ctx, in)
		if err != nil {
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.Error(err))
			return
		}

		res := &MessageCreateRes{out.Message.Id}
		ginctx.JSON(http.StatusCreated, res)
	}
}
