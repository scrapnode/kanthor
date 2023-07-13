package interceptors

import (
	gatewayinterceptors "github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	grpccore "google.golang.org/grpc"
)

func Interceptors() gatewayinterceptors.Interceptors {
	mw := gatewayinterceptors.Interceptors{
		Unary:  []grpccore.UnaryServerInterceptor{},
		Stream: []grpccore.StreamServerInterceptor{},
	}
	return mw
}
