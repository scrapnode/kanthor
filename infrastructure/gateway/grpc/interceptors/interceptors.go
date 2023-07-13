package interceptors

import grpccore "google.golang.org/grpc"

type Interceptors struct {
	Unary  []grpccore.UnaryServerInterceptor
	Stream []grpccore.StreamServerInterceptor
}
