package auth

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/metadata"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	grpccore "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func UnaryServerInterceptor(
	logger logging.Logger,
	engine authenticator.Authenticator,
) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		ctx, err = authenticate(logger, engine, ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(
	logger logging.Logger,
	engine authenticator.Authenticator,
) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		ctx, err := authenticate(logger, engine, ss.Context())
		if err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		wrapped := stream.WrapServerStream(ss)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

func authenticate(
	logger logging.Logger,
	engine authenticator.Authenticator,
	ctx context.Context,
) (context.Context, error) {
	t, err := token(ctx, engine.Scheme())
	if err != nil {
		logger.Error(err.Error())
		return ctx, err
	}

	account, err := engine.Verify(t)
	if err != nil {
		logger.Error(err.Error())
		return ctx, err
	}

	logger.Debugw("authenticated", "account_sub", account.Sub)
	return authenticator.AccountWithContext(ctx, account), nil
}

func token(ctx context.Context, scheme string) (string, error) {
	authorization := metadata.Metadata(ctx).Get("authorization")[0]
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
