package grpc

import (
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/auth"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/measurement"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/recovery"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	grpccore "google.golang.org/grpc"
)

func New(
	logger logging.Logger,
	meter metric.Meter,
	authenticator authenticator.Authenticator,
) *grpccore.Server {
	return grpccore.NewServer(
		grpccore.ChainUnaryInterceptor(
			measurement.UnaryServerInterceptor(meter),
			recovery.UnaryServerInterceptor(logger),
			auth.UnaryServerInterceptor(logger, authenticator),
		),
		grpccore.ChainStreamInterceptor(
			measurement.StreamServerInterceptor(meter),
			recovery.StreamServerInterceptor(logger),
			auth.StreamServerInterceptor(logger, authenticator),
		),
	)
}
