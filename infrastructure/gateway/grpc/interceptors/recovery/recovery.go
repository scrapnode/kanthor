package recovery

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	grpccore "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
)

func UnaryServerInterceptor(logger logging.Logger) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 64<<10)
				stack = stack[:runtime.Stack(stack, false)]

				logger.Errorw("SYSTEM.ERROR.RECOVER", "panic", r, "stack", stack)
				err = status.Error(codes.Internal, "SYSTEM.ERROR")
			}
		}()

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
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 64<<10)
				stack = stack[:runtime.Stack(stack, false)]

				logger.Errorw("SYSTEM.ERROR.RECOVER", "panic", r, "stack", stack)
				err = status.Error(codes.Internal, "SYSTEM.ERROR")
			}
		}()

		return handler(srv, ss)
	}
}
