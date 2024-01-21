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
				ginctx.AbortWithStatusJSON(http.StatusBadRequest, gateway.ErrorString("GATEWAY.MIDDLEWARE.IDEMPOTENCY_KEY.EMPTY.ERROR"))
				return
			}

			ctx := ginctx.MustGet(gateway.Ctx).(context.Context)
			ok, err := engine.Validate(ctx, key)
			if err != nil {
				ginctx.AbortWithStatusJSON(http.StatusInternalServerError, gateway.ErrorString("GATEWAY.MIDDLEWARE.IDEMPOTENCY_KEY.VERIFY.ERROR"))
				return
			}

			if !ok {
				ginctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gateway.ErrorString("GATEWAY.MIDDLEWARE.IDEMPOTENCY_KEY.INVALID.ERROR"))
				return
			}
		}

		ginctx.Next()
	}
}
