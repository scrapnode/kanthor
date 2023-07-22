package auth

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	grpccore "google.golang.org/grpc"
	"strings"
)

// DefaultPublic returns the list of public method
func DefaultPublic() map[string]bool {
	public := map[string]bool{
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
		"/grpc.health.v1.Health/Check":                                   true,
		"/grpc.health.v1.Health/Watch":                                   true,
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

		ctx, err = authenticate(logger, engine, ctx, public, method)
		if err != nil {
			return nil, err
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

		ctx, err := authenticate(logger, engine, ss.Context(), public, method)
		if err != nil {
			return err
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
	public map[string]bool,
	method string,
) (context.Context, error) {
	if publicable(method, public) {
		return context.WithValue(ctx, gateway.AccessPublicable, true), nil
	}

	t, err := TokenFromContext(ctx, engine.Scheme())
	if err != nil {
		logger.Errorw(err.Error())
		return ctx, err
	}

	account, err := engine.Verify(t)
	if err != nil {
		logger.Errorw(err.Error())
		return ctx, gateway.Err401("AUTH.VERIFY_FAILED")
	}

	logger.Debugw("authenticated", "account_sub", account.Sub)
	return context.WithValue(ctx, authenticator.CtxAuthAccount, account), nil
}

func TokenFromContext(ctx context.Context, scheme string) (string, error) {
	meta := gateway.ExtractIncoming(ctx)

	authorization := meta.Get("authorization")
	if authorization == "" {
		return "", gateway.Err401("AUTH.HEADERS.AUTHORIZATION_EMPTY")
	}

	segments := strings.Split(authorization, " ")
	if len(segments) != 2 {
		return "", gateway.Err401("AUTH.HEADERS.AUTHORIZATION_MALFORMED")
	}
	if !strings.EqualFold(segments[0], scheme) {
		return "", gateway.Err401("AUTH.HEADERS.AUTHORIZATION_SCHEME_INVALID")
	}

	return segments[1], nil
}
