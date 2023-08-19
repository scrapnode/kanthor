package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"time"
)

func UseMetrics(metrics metrics.Metrics) gin.HandlerFunc {
	dur := "http_request_duration_seconds"
	err4xx := "http_requests_4xx_error"
	err5xx := "http_requests_5xx_error"

	return func(ginctx *gin.Context) {
		start := time.Now()

		ginctx.Next()

		duration := time.Now().Sub(start)
		metrics.Observe(dur, duration.Seconds())

		status := ginctx.Writer.Status()
		if status >= 400 && status < 500 {
			metrics.Count(err4xx, 1)
		}
		if status >= 500 && status < 600 {
			metrics.Count(err5xx, 1)
		}
	}
}
