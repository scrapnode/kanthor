package recovery

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	grpccore "google.golang.org/grpc"
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

				logger.Errorw("recover from panic", "panic", r, "stack", string(stack))
				err = gateway.Err500("SYSTEM.INTERNAL")
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

				logger.Errorw("recover from panic", "panic", r, "stack", string(stack))
				err = gateway.Err500("SYSTEM.INTERNAL")
			}
		}()

		return handler(srv, ss)
	}
}
