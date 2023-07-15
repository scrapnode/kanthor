package grpc

import (
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
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

	return grpccore.NewServer(
		grpccore.ChainUnaryInterceptor(ics.Unary...),
		grpccore.ChainStreamInterceptor(ics.Stream...),
	)
}
