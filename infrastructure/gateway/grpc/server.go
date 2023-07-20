package grpc

import (
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/protos"
	grpccore "google.golang.org/grpc"
)

func NewServer(withInterceptors ...interceptors.InterceptorOps) *grpccore.Server {
	ics := interceptors.Interceptors{
		Unary:  []grpccore.UnaryServerInterceptor{},
		Stream: []grpccore.StreamServerInterceptor{},
	}
	for _, withInterceptor := range withInterceptors {
		withInterceptor(&ics)
	}

	server := grpccore.NewServer(
		grpccore.ChainUnaryInterceptor(ics.Unary...),
		grpccore.ChainStreamInterceptor(ics.Stream...),
	)
	protos.RegisterHealthServer(server, &healthcheck{})

	return server
}
