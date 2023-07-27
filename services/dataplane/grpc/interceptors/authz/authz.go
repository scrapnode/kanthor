package authz

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
	grpccore "google.golang.org/grpc"
	"strings"
)

func DefaultProtected() map[string]bool {
	protected := map[string]bool{}
	return protected
}

func UnaryServerInterceptor(
	logger logging.Logger,
	engine authorizator.Authorizator,
	protected map[string]bool,
) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		ctx, err = authorize(logger.With("method", info.FullMethod), engine, ctx, protected, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(
	logger logging.Logger,
	engine authorizator.Authorizator,
	protected map[string]bool,
) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		ctx, err := authorize(logger.With("method", info.FullMethod), engine, ss.Context(), protected, info.FullMethod)
		if err != nil {
			return err
		}

		wrapped := stream.WrapServerStream(ss)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

func authorize(
	logger logging.Logger,
	engine authorizator.Authorizator,
	ctx context.Context,
	protected map[string]bool,
	method string,
) (context.Context, error) {
	if public, ok := ctx.Value(gateway.AccessPublicable).(bool); ok && public {
		return ctx, nil
	}

	if protectable(method, protected) {
		return context.WithValue(ctx, gateway.AccessProtectable, true), nil
	}

	acc, ok := ctx.Value(authenticator.CtxAuthAccount).(*authenticator.Account)
	if !ok {
		return ctx, gateway.Err401("ACCOUNT.UNKNOWN")
	}

	app, ok := ctx.Value(usecase.CtxApplication).(*entities.Application)
	if !ok {
		return ctx, gateway.Err401("APPLICATION.UNKNOWN")
	}

	obj, act, err := action(method)
	if err != nil {
		logger.Error(err.Error())
		return ctx, gateway.Err400("METHOD.INVALID")
	}

	ok, err = engine.Enforce(app.Id, acc.Sub, obj, act)
	if err != nil {
		logger.Errorw(err.Error(), "app_id", app.Id, "account_sub", acc.Sub)

		return ctx, gateway.Err500("AUTHORIZATOR.INTERNAL")
	}

	if !ok {
		return ctx, gateway.Err403("AUTHORIZATOR.FORBIDDEN")
	}

	logger.Debugw("granted permission", "app_id", app.Id, "account_sub", acc.Sub)
	return ctx, nil
}

func action(method string) (string, string, error) {
	segments := strings.Split(method, "/")
	if len(segments) != 3 {
		return "", "", errors.New("malformed method")
	}

	return segments[1], segments[2], nil
}

func protectable(value string, maps map[string]bool) bool {
	should, ok := maps[value]
	return ok && should
}
