package interceptors

import (
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/auth"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/log"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/measurement"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/recovery"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	grpccore "google.golang.org/grpc"
)

type Interceptors struct {
	Unary  []grpccore.UnaryServerInterceptor
	Stream []grpccore.StreamServerInterceptor
}

type InterceptorOps func(ics *Interceptors)

func WithRecovery(logger logging.Logger) InterceptorOps {
	return func(ics *Interceptors) {
		ics.Unary = append(ics.Unary, recovery.UnaryServerInterceptor(logger))
		ics.Stream = append(ics.Stream, recovery.StreamServerInterceptor(logger))
	}
}

func WithLog(logger logging.Logger) InterceptorOps {
	return func(ics *Interceptors) {
		ics.Unary = append(ics.Unary, log.UnaryServerInterceptor(logger))
		ics.Stream = append(ics.Stream, log.StreamServerInterceptor(logger))
	}
}

func WithAuth(logger logging.Logger, engine authenticator.Authenticator, public map[string]bool) InterceptorOps {
	return func(ics *Interceptors) {
		ics.Unary = append(ics.Unary, auth.UnaryServerInterceptor(logger, engine, public))
		ics.Stream = append(ics.Stream, auth.StreamServerInterceptor(logger, engine, public))
	}
}

func WithMeasurement(meter metric.Meter) InterceptorOps {
	return func(ics *Interceptors) {
		ics.Unary = append(ics.Unary, measurement.UnaryServerInterceptor(meter))
		ics.Stream = append(ics.Stream, measurement.StreamServerInterceptor(meter))
	}
}
