package middlewares

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/telemetry"
	"go.opentelemetry.io/otel/attribute"
)

func UseTracing(name string) gin.HandlerFunc {
	tracer := telemetry.Tracer(name)

	return func(ginctx *gin.Context) {
		method := ginctx.Request.Method
		route := ginctx.FullPath()
		url := fmt.Sprintf("%s?%s", ginctx.Request.URL.Path, ginctx.Request.URL.RawQuery)

		gwwtx := ginctx.MustGet(gateway.Ctx).(context.Context)
		ctx, main := tracer.Start(gwwtx, fmt.Sprintf("%s %s", method, route))
		defer main.End()

		main.SetAttributes(
			attribute.String("http.method", method),
			attribute.String("http.route", route),
			attribute.String("http.url", url),
		)

		// override gateway context so inside the app we can use it to continue our tracing
		ginctx.Set(gateway.Ctx, ctx)
		// Process request
		ginctx.Next()

		main.SetAttributes(
			attribute.Int("http.status_code", ginctx.Writer.Status()),
		)
	}
}
