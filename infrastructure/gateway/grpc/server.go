package grpc

import (
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/auth"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/measurement"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/recovery"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	grpccore "google.golang.org/grpc"
)

func NewServer(
	logger logging.Logger,
	meter metric.Meter,
	authenticator authenticator.Authenticator,
	interceptors interceptors.Interceptors,
) *grpccore.Server {
	logger = logger.With("gateway", "grpc")

	unary := append(
		[]grpccore.UnaryServerInterceptor{
			measurement.UnaryServerInterceptor(meter),
			recovery.UnaryServerInterceptor(logger),
			auth.UnaryServerInterceptor(logger, authenticator),
		},
		interceptors.Unary...,
	)

	stream := append(
		[]grpccore.StreamServerInterceptor{
			measurement.StreamServerInterceptor(meter),
			recovery.StreamServerInterceptor(logger),
			auth.StreamServerInterceptor(logger, authenticator),
		},
		interceptors.Stream...,
	)
	return grpccore.NewServer(
		grpccore.ChainUnaryInterceptor(unary...),
		grpccore.ChainStreamInterceptor(stream...),
	)
}
