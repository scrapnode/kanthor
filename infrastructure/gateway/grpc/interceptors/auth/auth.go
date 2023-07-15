package auth

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	grpccore "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func DefaultPublic() map[string]bool {
	public := map[string]bool{
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
	}

	return public
}

func UnaryServerInterceptor(
	logger logging.Logger,
	engine authenticator.Authenticator,
	public map[string]bool,
) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		method := info.FullMethod

		if publicable(method, public) {
			ctx = context.WithValue(ctx, gateway.AccessPublicable, true)
			return handler(ctx, req)
		}

		ctx, err = authenticate(logger, engine, ctx)
		if err != nil {
			logger.Errorw(err.Error(), "method", method)
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(
	logger logging.Logger,
	engine authenticator.Authenticator,
	public map[string]bool,
) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		method := info.FullMethod
		wrapped := stream.WrapServerStream(ss)

		if publicable(method, public) {
			wrapped.WrappedContext = context.WithValue(wrapped.WrappedContext, gateway.AccessPublicable, true)
			return handler(srv, wrapped)
		}

		ctx, err := authenticate(logger, engine, ss.Context())
		if err != nil {
			logger.Errorw(err.Error(), "method", method)
			return status.Error(codes.Unauthenticated, err.Error())
		}
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

func publicable(value string, maps map[string]bool) bool {
	should, ok := maps[value]
	return ok && should
}

func authenticate(
	logger logging.Logger,
	engine authenticator.Authenticator,
	ctx context.Context,
) (context.Context, error) {
	t, err := token(ctx, engine.Scheme())
	if err != nil {
		return ctx, err
	}

	account, err := engine.Verify(t)
	if err != nil {
		return ctx, err
	}

	logger.Debugw("authenticated", "account_sub", account.Sub)
	return context.WithValue(ctx, authenticator.CtxAuthAccount, account), nil
}

func token(ctx context.Context, scheme string) (string, error) {
	meta := gateway.ExtractIncoming(ctx)

	authorization := meta.Get("authorization")
	if authorization == "" {
		return "", status.Error(codes.Unauthenticated, "AUTH.HEADERS.AUTHORIZATION_EMPTY")
	}

	segments := strings.Split(authorization, " ")
	if len(segments) != 2 {
		return "", status.Errorf(codes.Unauthenticated, "AUTH.HEADERS.AUTHORIZATION_MALFORMED")
	}
	if !strings.EqualFold(segments[0], scheme) {
		return "", status.Errorf(codes.Unauthenticated, "AUTH.HEADERS.AUTHORIZATION_SCHEME_INVALID")
	}

	return segments[1], nil
}
