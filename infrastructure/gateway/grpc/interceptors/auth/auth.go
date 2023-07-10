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

func UnaryServerInterceptor(logger logging.Logger, engine authenticator.Authenticator) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		t, err := token(ctx, engine.Scheme())
		if err != nil {
			return nil, err
		}

		account, err := engine.Verify(t)
		if err != nil {
			return nil, err
		}

		logger.Debugw("authenticated", "account_sub", account.Sub)
		ctx = context.WithValue(ctx, authenticator.CtxAuthAccount, account)
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(logger logging.Logger, engine authenticator.Authenticator) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		// @TODO: log some useful information
		ctx := ss.Context()
		t, err := token(ctx, engine.Scheme())
		if err != nil {
			return err
		}

		account, err := engine.Verify(t)
		if err != nil {
			return err
		}

		wrapped := stream.WrapServerStream(ss)
		wrapped.WrappedContext = context.WithValue(ctx, authenticator.CtxAuthAccount, account)
		return handler(srv, wrapped)
	}
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
