package measurement

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	grpccore "google.golang.org/grpc"
	"time"
)

func UnaryServerInterceptor(
	meter metric.Meter,
) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		defer func() {
			duration := time.Now().Sub(start).Seconds()
			meter.Histogram(
				gateway.MetricLatency, duration,
				metric.Label("grpc_type", "unary"),
			)
		}()

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(
	meter metric.Meter,
) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		start := time.Now()
		defer func() {
			duration := time.Now().Sub(start).Seconds()
			meter.Histogram(
				gateway.MetricLatency, duration,
				metric.Label("step", "token_resolve"),
				metric.Label("grpc_type", "stream"),
			)
		}()

		return handler(srv, ss)
	}
}
