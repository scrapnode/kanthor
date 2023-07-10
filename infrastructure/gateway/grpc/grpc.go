package grpc

import (
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/auth"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	grpccore "google.golang.org/grpc"
)

func New(
	logger logging.Logger,
	authenticator authenticator.Authenticator,
) *grpccore.Server {
	return grpccore.NewServer(
		grpccore.ChainUnaryInterceptor(auth.UnaryServerInterceptor(logger, authenticator)),
		grpccore.ChainStreamInterceptor(auth.StreamServerInterceptor(logger, authenticator)),
	)
}
