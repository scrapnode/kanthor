package log

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	grpccore "google.golang.org/grpc"
)

func UnaryServerInterceptor(logger logging.Logger) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		logger.Debugw(info.FullMethod, "request", utils.Stringify(req))
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(logger logging.Logger) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) (err error) {
		logger.Debugw(info.FullMethod)
		return handler(srv, ss)
	}
}
