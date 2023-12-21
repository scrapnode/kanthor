package middlewares

import (
	"context"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/logging"
)

var HeaderIdempotencyKey = "idempotency-key"
var requires = []string{
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

func UseIdempotency(logger logging.Logger, engine idempotency.Idempotency, bypass bool) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		method := ginctx.Request.Method
		required := slices.Contains(requires, method) && !bypass

		if required {
			key := ginctx.GetHeader(HeaderIdempotencyKey)
			if key == "" {
				ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.NewError("no idempotency key"))
				return
			}

			ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
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
		}

		ginctx.Next()
	}
}
