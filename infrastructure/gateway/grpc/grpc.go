package grpc

import (
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/auth"
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
		grpccore.ChainUnaryInterceptor(auth.UnaryServerInterceptor(logger, meter, authenticator)),
		grpccore.ChainStreamInterceptor(auth.StreamServerInterceptor(logger, meter, authenticator)),
	)
}
