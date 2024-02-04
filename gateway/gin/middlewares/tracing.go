package middlewares

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func UseTracing(name string) gin.HandlerFunc {
	tracer := telemetry.Tracer(name)

	return func(ginctx *gin.Context) {
		method := ginctx.Request.Method
		route := ginctx.FullPath()
		url := fmt.Sprintf("%s?%s", ginctx.Request.URL.Path, ginctx.Request.URL.RawQuery)

		gwwtx := context.WithValue(ginctx.MustGet(gateway.Ctx).(context.Context), telemetry.CtxTracer, tracer)
		attributes := trace.WithAttributes(
			attribute.String("http.method", method),
			attribute.String("http.route", route),
			attribute.String("http.url", url),
		)
		ctx, span := tracer.Start(gwwtx, fmt.Sprintf("%s %s", method, route), attributes)
		defer span.End()

		// override gateway context so inside the app we can use it to continue our tracing
		ginctx.Set(gateway.Ctx, ctx)
		// Process request
		ginctx.Next()

		span.SetAttributes(
			attribute.Int("http.status_code", ginctx.Writer.Status()),
		)
	}
}
