package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

var HeaderIdempotencyKey = "idempotency-key"

func UseIdempotency(logger logging.Logger, engine idempotency.Idempotency) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		method := ginctx.Request.Method
		bypass := method == http.MethodGet || method == http.MethodHead
		if bypass {
			logger.Debugw("bypass method", "method", method)
			ginctx.Next()
			return
		}

		key := ginctx.GetHeader(HeaderIdempotencyKey)
		if key == "" {
			logger.Warnw("no idempotency key", "method", method)
			ginctx.Next()
			return
		}

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		ok, err := engine.Validate(ctx, key)
		if err != nil {
			logger.Error(err)
			ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.NewError("could not validate idempotency key"))
			return
		}

		if !ok {
			logger.Errorw("duplicated request", "method", method, "idempotency_key", key)
			ginctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gateway.NewError("duplicated request"))
			return
		}

		ginctx.Next()
	}
}
