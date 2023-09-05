package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"time"
)

func UseMetrics(metrics metric.Metrics) gin.HandlerFunc {
	dur := "http_request_duration_seconds"
	err4xx := "http_requests_4xx_error"
	err5xx := "http_requests_5xx_error"

	return func(ginctx *gin.Context) {
		start := time.Now()

		ginctx.Next()

		ctx := ginctx.MustGet(gateway.KeyCtx).(context.Context)
		duration := time.Now().Sub(start)
		metrics.Observe(ctx, dur, duration.Seconds())

		status := ginctx.Writer.Status()
		if status >= 400 && status < 500 {
			metrics.Count(ctx, err4xx, 1)
		}
		if status >= 500 && status < 600 {
			metrics.Count(ctx, err5xx, 1)
		}
	}
}
