package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
)

func UseMetric(metrics metric.Metric, service string) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		start := time.Now()

		ginctx.Next()

		ctx := ginctx.MustGet(gateway.KeyContext).(context.Context)
		duration := time.Since(start)
		metrics.Observe(ctx, service, metric.HttpRequestDuration, duration.Seconds())

		status := ginctx.Writer.Status()
		if status >= 400 && status < 500 {
			metrics.Count(ctx, service, metric.HttpRequestErr4xxx, 1)
		}
		if status >= 500 && status < 600 {
			metrics.Count(ctx, service, metric.HttpRequestErr5xxx, 1)
		}
	}
}
