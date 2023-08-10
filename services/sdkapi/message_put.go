package sdkapi

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/validator"
	usecase "github.com/scrapnode/kanthor/usecases/sdk"
	"net/http"
)

type MessagePutReq struct {
	Type string `json:"type" binding:"required"`

	Body     map[string]interface{} `json:"body" binding:"required"`
	Headers  map[string]string      `json:"headers"`
	Metadata map[string]string      `json:"metadata"`
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
func UseMessagePut(logger logging.Logger, validator validator.Validator, uc usecase.Sdk) gin.HandlerFunc {
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

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		appId := ginctx.Param("app_id")
		headers := http.Header{}
		if len(req.Headers) > 0 {
			for k, v := range req.Headers {
				headers.Set(k, v)
			}
		}
		ucreq := &usecase.MessagePutReq{
			AppId:    appId,
			Type:     req.Type,
			Body:     body,
			Headers:  headers,
			Metadata: req.Metadata,
		}
		if err := validator.Struct(ucreq); err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("invalid request"))
			return
		}

		ucres, err := uc.Message().Put(ctx, ucreq)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("oops, something went wrong"))
			return
		}

		res := &MessagePutRes{ucres.Msg.Id}
		ginctx.JSON(http.StatusCreated, res)
	}
}
